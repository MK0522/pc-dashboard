let systemInfo = null
let systemInfoData = null
let isConnected = true
let connectionStatusElement = null
let connectionMessageElement = null
const FETCH_TIMEOUT_MS = 1500
let lastSuccessfulPollAt = Date.now()
let connectionWatchdogId = null

async function fetchJsonWithTimeout(url) {
	const controller = new AbortController()
	const timeoutId = setTimeout(() => controller.abort(), FETCH_TIMEOUT_MS)

	try {
		const res = await fetch(url, { signal: controller.signal })
		if (!res.ok) {
			throw new Error(`HTTP ${res.status}`)
		}

		return await res.json()
	} finally {
		clearTimeout(timeoutId)
	}
}

// 연결 상태 표시
function setConnectionStatus(connected, message = null) {
	if (!connectionStatusElement || !connectionMessageElement) {
		return
	}

	isConnected = connected

	if (connected) {
		connectionStatusElement.classList.add("reconnected")
		connectionMessageElement.innerText = "✓ 연결되었습니다"
		connectionStatusElement.classList.remove("hidden")

		// 2초 후 사라지기
		setTimeout(() => {
			connectionStatusElement.classList.add("hidden")
			connectionStatusElement.classList.remove("reconnected")
		}, 2000)
	} else {
		connectionStatusElement.classList.remove("reconnected")
		connectionMessageElement.innerText = message || "⚠️ 연결이 끊겼습니다"
		connectionStatusElement.classList.remove("hidden")
	}
}

// 시스템 고정 정보 로딩
async function loadSystemInfo() {
	try {
		const data = await fetchJsonWithTimeout("/api/info")

    systemInfo = data

		osName.innerText = data.osName
		osVersion.innerText = data.osVersion

		document.getElementById("cpu-name").innerText =
			data.cpuName

		document.getElementById("cpu-cores").innerText =
			`Cores: ${data.cpuCores}`

		document.getElementById("cpu-threads").innerText =
			`Threads: ${data.cpuThreads}`

		document.getElementById("ram-total").innerText =
			`${data.ramTotal.toFixed(1)} GB`

		systemInfoData = data

		if (!isConnected) {
			setConnectionStatus(true)
		}
		lastSuccessfulPollAt = Date.now()
	} catch (err) {
		console.error("시스템 정보 로드 실패:", err)
		if (isConnected) {
			setConnectionStatus(false, "⚠️ 서버가 응답하지 않습니다")
		}
	}
}

const osName =	document.getElementById("os-name")
const osVersion =	document.getElementById("os-version")
const uptimeText =	document.getElementById("uptime")

const cpuText = document.getElementById("cpu")
const cpuBar = document.getElementById("cpu-bar")
const cpuGrid =	document.getElementById("cpu-grid")
const toggleButton =	document.getElementById("toggle-core-button")

const ramText = document.getElementById("ram")
const ramBar = document.getElementById("ram-bar")
const ramUsedText = document.getElementById("ram-used")

const diskContainer =	document.getElementById("disk-container")
const networkContainer =	document.getElementById("network-container")

// 실시간 상태 업데이트
async function updateStats() {
	try {
		const data = await fetchJsonWithTimeout("/api/stats")

		uptimeText.innerText =	formatUptime(data.hostUptime)

		// CPU 사용률
		cpuText.innerText = data.cpuUsage.toFixed(2)
		cpuBar.style.width = data.cpuUsage + "%"
		updateBarColor(cpuBar, data.cpuUsage)
		renderPerCPU(data.perCpuUsage)

		// RAM 사용률
		ramText.innerText = data.ramUsage.toFixed(2)
		ramBar.style.width = data.ramUsage + "%"
		updateBarColor(ramBar, data.ramUsage)
		ramUsedText.innerText =	`${data.ramUsed.toFixed(1)} GB`

		renderDisks(
			systemInfoData?.disks ?? [],
			data.disks
		)

		renderNetworks(
			systemInfoData?.networks ?? [],
			data.networks ?? []
		)

		if (!isConnected) {
			setConnectionStatus(true)
		}
		lastSuccessfulPollAt = Date.now()
	} catch (err) {
		console.error("통계 업데이트 실패:", err)
		if (isConnected) {
			setConnectionStatus(false, "⚠️ 데이터를 가져올 수 없습니다")
		}
	}
}

function startConnectionWatchdog() {
	if (connectionWatchdogId !== null) {
		clearInterval(connectionWatchdogId)
	}

	connectionWatchdogId = setInterval(() => {
		const elapsed = Date.now() - lastSuccessfulPollAt

		if (elapsed > 3500 && isConnected) {
			setConnectionStatus(false, "⚠️ 연결이 끊겼습니다")
		}
	}, 500)
}

toggleButton.addEventListener("click", () => {

	cpuGrid.classList.toggle("hidden")

	if (cpuGrid.classList.contains("hidden")) {
		toggleButton.innerText =
			"코어별 사용률 보기"
	}
	else {
		toggleButton.innerText =
			"숨기기"
	}
})

function formatUptime(seconds) {

	const days =
		Math.floor(seconds / 86400)

	const hours =
		Math.floor((seconds % 86400) / 3600)

	const minutes =
		Math.floor((seconds % 3600) / 60)

	const secs =
	seconds % 60

	let result = "Up "

	if (days > 0) {
		result += `${days}d `
	}

	if (hours > 0 || days > 0) {
		result += `${hours}h `
	}

	result += `${minutes}m`
	result += ` ${secs}s`

	return result
}

function renderPerCPU(perCpuUsage) {

	cpuGrid.innerHTML = ""

	perCpuUsage.forEach((usage, index) => {

		const core = document.createElement("div")
		core.className = "cpu-core"

		// 색상 계산
		const hue = 120 - (usage * 1.2)

		core.style.backgroundColor =
			`hsl(${hue}, 70%, 25%)`

		core.innerHTML = `
			<div class="core-name">
				CPU-${index}
			</div>

			<div class="core-usage">
				${usage.toFixed(0)}%
			</div>
		`

		cpuGrid.appendChild(core)
	})
}

function renderDisks(infoDisks, statDisks) {

	diskContainer.innerHTML = ""

	infoDisks.forEach((infoDisk, index) => {

		const statDisk = statDisks[index]

		if (!statDisk) {
			return
		}

		const usage =
			statDisk.usagePercent

		const hue =
			120 - (usage * 1.2)

		const diskCard =
			document.createElement("div")

		diskCard.className =
			"disk-item"

		diskCard.innerHTML = `
			<div class="disk-top">
				<div class="disk-name">
					${infoDisk.name}
				</div>

				<div class="disk-usage">
					${usage.toFixed(0)}%
				</div>
			</div>

			<div class="disk-info">
				<div class="disk-capacity">
					${statDisk.usedGB.toFixed(0)} GB
					/
					${infoDisk.totalGB.toFixed(0)} GB
				</div>
				<div class="disk-speed">
					↓ ${statDisk.readMBs.toFixed(1)} MB/s
					&nbsp;
					↑ ${statDisk.writeMBs.toFixed(1)} MB/s
				</div>
			</div>

			<div class="bar">
				<div
					class="fill"
					style="
						width: ${usage}%;
						background-color:
							hsl(${hue}, 80%, 50%);
					"
				></div>
			</div>
		`

		diskContainer.appendChild(diskCard)
	})
}

function renderNetworks(infoNetworks, statNetworks) {

	networkContainer.innerHTML = ""

	const statMap = new Map(
		statNetworks.map((network) => [network.name, network])
	)

	let visibleCount = 0

	infoNetworks.forEach((infoNetwork) => {

		const statNetwork = statMap.get(infoNetwork.name)

		if (!statNetwork) {
			return
		}

		const isIdle =
			statNetwork.downloadMBs < 0.01 &&
			statNetwork.uploadMBs < 0.01

		if (isIdle && isLikelyVirtualInterface(infoNetwork.name)) {
			return
		}

		visibleCount++

		const networkCard =
			document.createElement("div")

		networkCard.className =
			"network-item"

		const downloadPercent =
			Math.min((statNetwork.downloadMBs / 10) * 100, 100)

		const uploadPercent =
			Math.min((statNetwork.uploadMBs / 10) * 100, 100)

		networkCard.innerHTML = `
			<div class="network-top">
				<div class="network-name">
					${infoNetwork.name}
				</div>
				<div class="network-badge">Ethernet</div>
			</div>

			<div class="network-metrics">
				<div class="network-row">
					<div class="network-label">Download</div>
					<div class="network-value down-value">
						${formatNetworkMBs(statNetwork.downloadMBs)} MB/s
					</div>
				</div>
				<div class="network-meter">
					<div
						class="network-fill down"
						style="width: ${downloadPercent}%;"
					></div>
				</div>

				<div class="network-row">
					<div class="network-label">Upload</div>
					<div class="network-value up-value">
						${formatNetworkMBs(statNetwork.uploadMBs)} MB/s
					</div>
				</div>
				<div class="network-meter">
					<div
						class="network-fill up"
						style="width: ${uploadPercent}%;"
					></div>
				</div>
			</div>
		`

		networkContainer.appendChild(networkCard)
	})

	if (visibleCount === 0) {
		networkContainer.innerHTML = `
			<div class="network-empty">
				활성 Ethernet 인터페이스를 찾는 중입니다.
			</div>
		`
	}
}

function isLikelyVirtualInterface(name) {
	const lower = name.toLowerCase()

	return (
		lower.includes("docker") ||
		lower.includes("veth") ||
		lower.includes("virbr") ||
		lower.includes("vbox") ||
		lower.includes("vmnet") ||
		lower.includes("vethernet") ||
		lower.includes("loopback") ||
		lower.startsWith("br-") ||
		lower.startsWith("tun") ||
		lower.startsWith("tap")
	)
}

function formatNetworkMBs(value) {
	if (value >= 10) {
		return value.toFixed(1)
	}

	if (value >= 1) {
		return value.toFixed(2)
	}

	return value.toFixed(3)
}

function updateBarColor(element, value) {
    hue = 120 - (value * 1.2)
	element.style.backgroundColor = `hsl(${hue}, 80%, 50%)`
}

// 최초 실행
(async () => {
	connectionStatusElement = document.getElementById("connection-status")
	connectionMessageElement = document.getElementById("connection-message")

	startConnectionWatchdog()

	await loadSystemInfo()
	await updateStats()

	// 1초마다 갱신
	setInterval(updateStats, 1000)
})()
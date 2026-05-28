# PC Dashboard

로컬 PC의 시스템 상태를 실시간으로 확인할 수 있는 경량 시스템 모니터링 대시보드입니다.

Go 백엔드와 Vanilla JavaScript 프론트엔드 기반으로 제작되었으며,
CPU / RAM / Storage / Network 상태를 웹 UI에서 실시간으로 확인할 수 있습니다.

![preview](./docs/preview.png)

---

# Features

## System Information
- 운영체제 정보
- CPU 모델 / 코어 / 스레드
- RAM 총 용량
- 디스크 목록 및 총 용량
- 활성 네트워크 인터페이스 감지

## Real-Time Monitoring
- CPU 전체 사용률
- 코어별 CPU 사용률
- RAM 사용률
- 디스크 사용량
- 디스크 읽기 / 쓰기 속도
- 네트워크 업로드 / 다운로드 속도
- 시스템 업타임

## Connection Monitoring
- Fetch timeout 기반 연결 감시
- Watchdog 기반 응답 상태 추적
- 연결 끊김 / 재연결 상태 UI 표시
- 브라우저 캐시 방지 처리

---

# Tech Stack

## Backend
- Go
- net/http
- gopsutil

## Frontend
- HTML5
- CSS3
- Vanilla JavaScript (ES6)

---

# Project Structure

```text
project/
├─ main.go
├─ api/
│  └─ handlers.go
├─ system/
│  ├─ info.go
│  ├─ network.go
│  ├─ stats.go
│  └─ types.go
├─ static/
│  ├─ index.html
│  ├─ style.css
│  ├─ app.js
│  └─ favicon.ico
```

---

# API

## `/api/info`
정적 시스템 정보 반환

### Includes
- OS 정보
- CPU 정보
- RAM 총량
- 디스크 정보
- 네트워크 인터페이스

---

## `/api/stats`
실시간 시스템 상태 반환

### Includes
- CPU 사용률
- 코어별 CPU 사용률
- RAM 사용률
- 디스크 사용량
- 디스크 I/O
- 네트워크 트래픽
- 시스템 업타임

---

# Run

## Run directly

```bash
go run .
```

---

## Build

```bash
go build -o pc-dashboard.exe
```

실행 후:

```text
http://localhost:8080
```

접속

---

# Design Goals

이 프로젝트는 다음 목표를 기준으로 개발되었습니다.

- 구현 우선
- 빠른 기능 검증
- 최소한의 구조
- 로컬 실행 중심
- 과한 추상화 지양

---

# Current Scope

현재 프로젝트는 다음 범위까지만 다룹니다.

- 단일 로컬 머신 모니터링
- 브라우저 기반 UI
- Windows 중심 개발 환경
- 로컬 네트워크 기준 사용

다음 항목은 현재 범위에서 제외했습니다.

- Docker 배포
- 인증 시스템
- GPU 모니터링
- 멀티 서버 관리
- Kubernetes / Redis / Database 연동

---

# Future Plans

- UI 테마 개선
- Electron 기반 데스크탑 앱 전환
- 시스템 알림 기능
- 프로세스 모니터링
- 더 안정적인 Cross Platform 지원

---

# Notes

- 브라우저 캐시 문제 방지를 위해 API 응답에 no-store 헤더를 적용했습니다.
- 연결이 끊기면 중앙 알림 UI가 표시됩니다.
- Watchdog 기반 상태 감시를 사용합니다.

---

# License

MIT
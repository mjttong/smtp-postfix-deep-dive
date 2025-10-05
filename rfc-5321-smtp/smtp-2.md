# 2. The SMTP Model

## 2.1. Basic Structure

User/File System ↔ Client SMTP ↔ Server SMTP ↔ File System

- SMTP 클라이언트의 역할
  - 전송할 메시지가 있는 경우 SMTP 서버에 양방향 전송 채널을 설정
  - 메일 메시지를 하나 이상의 SMTP 서버에 전송 및 전송 실패 보고
- SMTP 클라이언트가 메일을 보낼 때, 목적지 도메인이 최종 목적지일 수도, 중계되는 중간 목적지일 수도 있음

## 2.2. The Extension Model

### 2.2.1. Background

- 확장 모델이란 클라이언트, 서버가 원래의 SMTP 요구사항을 넘어서는 공유 기능을 활용하는 것을 말함
- 현재의 SMTP 구현은 기본 확장 매커니즘을 반드시 지원해야 함
  - 서버는 `EHLO`를 반드시 지원, 클라이언트는 `EHLO`를 우선 사용

## 2.3. SMTP Terminology

### 2.3.1. Mail Objects

- SMTP는 SMTP 봉투와 컨텐츠를 포함하는 메일 객체를 전송
- SMTP 봉투는 일련의 SMTP 프로토콜 단위로 구성
  - 발신자 주소, 수신자 주소, 선택적 프로토콜 확장 자료로 구성
- SMTP 컨텐츠는 SMTP DATA 프로토콜 단위로 구성되며 헤더와 본문으로 나뉨
  - 헤더는 헤더 이름, 콜론, 데이터 필드의 모음으로 구성, 메시지 형식 지정(RFC 5322)에 따라 구조화

### 2.3.10. Originator, Delivery, Relay, and Gateway Systems

- 수행하는 역할에 따른 SMTP 시스템을 구분함
- originating system(SMTP originator)
  - 메일을 외부 환경으로 내보내는 첫 지점
- delivery SMTP system
  - 메일을 받아 메시지 저장소에 저장
- relay SMTP system
  - SMTP 클라이언트로부터 메일을 받아 다른 SMTP 서버로 전송
  - 메일 메시지는 추적 정보 추가 외에는 변형하지 않음
- gateway SMTP system
  - 전송 환경의 클라이언트 시스템으로부터 메일을 받아 다른 전송 환경의 서버 시스템으로 전송
  - 릴레이와 유사하지만, 메시지 변환을 수행할 수 있다는 차이점 존재

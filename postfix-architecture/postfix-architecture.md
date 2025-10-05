# Postfix Architecture

🔗 [Postfix Architecture Overview](https://www.postfix.org/OVERVIEW.html)

![Image](/postfix-architecture/image/postfix-architecture.png)

## How Postfix receives mail

![Image](/postfix-architecture/image/receive.png)

### 메일 소스

- 해당 소스는 모두 `cleanup`으로 이동
- 네트워크 메일
  - `smtpd`, `qmqpd` 서버를 통해 네트워크 메일 유입
  - 메일 프로토콜의 캡슐화 제거, 메일 무결성 검사 수행
- 로컬 메일 submission
  - `sendmail` 명령 등을 통해 로컬 메일 제출
  - 해당 메일은 `postdrop` 명령에 의해 `maildrop` 큐에 저장
  - `pickup` 서버가 큐에 있는 로컬 메일을 가져와 무결성 검사 수행
- 내부 소스 메일
  - 아래의 메일이 내부 소스 메일에 해당
    - `local` 전송 에이전트가 전달한 메일
    - `bounce` 서버가 발신자에게 반송하는 메일
    - postmaster 알림

### `cleanup` 서버

- 메일이 큐에 저장되기 전 최종 처리를 진행
  - 누락된 `From:` 헤더, 기타 메시지 헤더 추가, 주소 변환
  - 정규표현식을 통한 컨텐츠 검사 수행
- 최종 처리 후 메일을 단일 파일로 incoming 큐에 저장
- 새로운 메일이 도착했음을 큐 관리자에게 알림

## How Postfix delivers mail

![Image](/postfix-architecture/image/deliver.png)

- 메일을 수신한 후 위의 과정을 거쳐 메일이 incoming 큐에 도착한 상황
- 해당 과정에서는 메일을 최종 목적지로 전달(delivers)

### `qmqr`

- 큐 관리자 서버로 메일을 누가 어떻게 보낼지를 결정
- incoming 큐와 deferred 큐에서 메일을 가져와 active 큐로 이동
- active 큐의 메일을 `smtp`, `lmtp`, `local`, `virtual`, `pipe` 등의 다양한 전달 에이전트에게 전달
- 만약 메일을 즉시 전달할 수 없는 경우 해당 메일은 deferred 큐에 따로 보관

### 전달 에이전트

- `smtp`: 네트워크를 통해 다른 SMTP 서버로 메일 전달
- `lmtp`: `smtp`와 유사, 다른 서버로 메일 전달
- `virtual`: UNIX 스타일 메일박스, maildir에만 메일 전달
- `local`: UNIX 스타일 메일 박스, maildir, aliases, .forward 파일 등을 처리 가능
  - `local`은 훅을 가지고 있어 메일박스 전달을 외부 명령이나 다른 postfix 전달 에이전트에게 위임하는 설정이 가능
  - 사용자 정의 스크립트, 프로그램에 메일을 전달하고자 할 때 중요
- `pipe`
  - 다른 메일 처리 시스템으로 메일을 보냄
  - 메일을 외부 스크립트, 프로그램으로 파이프함

# postfix 기초 설정

## Basic Configuration

- `myhostname`
  - postfix가 실행되는 서버의 FDQN을 지정
- `myorigin`
  - 발신 메일의 도메인 부분을 결정(`MAIL FROM`의 `reverse-path`)
  - `myhostname` 사용: 메일이 특정 호스트에서 발송된 것처럼 보임
  - `mydomain` 사용: 메일이 특정 도메인에서 온 것처럼 보임
  - 사용자가 불완전한 주소(user, root)로 메일을 보낼 때, Posifix가 자동으로 `myorigin`의 값을 붙여 완전한 주소로 만들어 메일 전송
- `mydestination`
  - Postfix가 로컬 배달의 최종 목적지라고 판단하는 도메인 목록
  - 이 목록에 포함된 도메인으로 수신된 메일은 다른 서버로 릴레이하지 않고 로컬에서 처리
  - 로컬 외 다른 서버의 도메인 이름을 쓰는 것이 아님, 해당 파라미터에는 로컬, 자신이 소유한 도메인 값만 써야 함
- `mynetworks`
  - 해당 서버를 릴레이 서버로 사용할 수 있는 클라이언트 정의
- `relay_domains`
  - 신뢰할 수 없는 클라이언트(`mynetworks` 외)로부터 받은 메일을 어떤 도메인으로 전달할지 결정
- `relayhost`
  - 메일을 직접 외부로 전송할지, 릴레이 호스트를 통해 전송할지 설정
- `inet_interfaces`
  - postfix 시스템이 수신 대기할 네트워크 인터페이스 지정
  - all: 모든 인터페이스에서 수신
  - localhost: 로컬 전용

## `aliases` 설정

```plain text
name: value1, value2, ...
```

- `name`: 수신자명
- `value`: 전달 대상(이메일 주소, 파일, 명령 등)

```plain text
# /etc/postfix/aliases:

postmaster: admin@example.com
root: postmaster
```

- `postmaster`, `root` 별칭 설정 필수
  - Postfix는 시스템 문제 발생 시 `postmaster`에게 알림 전송
  - `postmaster`는 실제 메일을 받을 수 있는 주소로 설정
- `RCPT TO` 명령의 수신자 주소가 로컬 도메인인 경우 `aliases` 검색
  - 일치하는 `name`이 있으면 `value`에 설정된 대상으로 메일 전달

## mailbox

- 이메일 메시지가 저장되는 물리적이거나 논리적인 공간
  - 물리적인 디렉터리, 파일일 수도 있으며 논리적인 개념으로만 존재할 수도 있음, 이는 mailbox의 형식에 따라 다름
- `Maildir` 형식(디렉터리-메시지 파일), `mbox`(단일 파일) 형식을 널리 사용
  - postfix의 `main.cf` 파일에서 어떤 형식의 mailbox를 사용할 지 결정할 수 있음

### 메일 배달 흐름

- `RCPT TO` 명령에서 수신자 주소를 기반으로 linux 사용자 계정을 찾음
- 해당 사용자의 mailbox에 메일을 저장
- 현대 메일 시스템에서는 리눅스 사용자 계정을 설정하는 대신 데이터베이스에 가상 사용자를 저장하고, `Maildir` 형식을 사용해 메일 저장

## Postfix 설정 예시

🔗 [Postfix Standard Configuration Examples: Postfix on a local network](https://www.postfix.org/STANDARD_CONFIGURATION_README.html#local_network)

- 하나의 메인 서버와 여러 대의 SMTP 서버가 서로 이메일을 주고받을 수 있는 로컬 네트워크 구조 형성
- 모든 서버
  - `user@example.com`으로 이메일 전송 가능
  - 자기 자신의 이름으로 오는 메일을 받을 수 있음
- 메인 서버
  - `user@example.com`로 오는 모든 메일을 최종적으로 받음
  - `mailhost.example.com`이라는 FDQN으로 설정되어 있음

### 일반 서버

```plain text
1 /etc/postfix/main.cf:
2 myorigin = $mydomain
3 mynetworks = 127.0.0.0/8 10.0.0.0/24
4 relay_domains =
5 # Optional: forward all non-local mail to mailhost
6 #relayhost = $mydomain
```

- 기능
  - 메일을 `user@example.com`로 보냄
  - 자기 자신에게 오는 메일 `user@hostname.example.com`을 최종적으로 받음
- 2행
  - 발신 메일 주소를 `user@example.com`으로 통일
  - 단, 해당 설정으로 인해 root, postmaster 등의 시스템 계정의 메일도 메인 서버로 전송됨
- 3행
  - 신뢰할 수 있는(릴레이 서버로 설정할 수 있는) 서버는 자기 자신 `127.0.0.0/8`과 자신의 로컬 네트워크 `10.0.0.0/24`
- 4행
  - 신뢰할 수 있는 서버 `mynetworks` 외의 서버는 릴레이 서버로 사용하지 않음
- 6행
  - 메인 서버를 발신 릴레이로 설정 시 해당 라인의 주석을 해제

### 메인 서버

```plain text
1 DNS:
2 example.com IN MX 10 mailhost.example.com.
3
4 /etc/postfix/main.cf:
5 myorigin = $mydomain
6 mydestination = $myhostname localhost.$mydomain localhost $mydomain
7 mynetworks = 127.0.0.0/8 10.0.0.0/24
8 relay_domains =
9 # Optional: forward all non-local mail to firewall
10 #relayhost = [firewall.example.com]
```

- 기능
  - 메일을 `user@example.com`로 보냄
  - 자기 자신에게 오는 메일 `user@hostname.example.com`을 최종적으로 받음
  - `user@example.com`으로 오는 메일을 최종적으로 받음
- DNS
  - MX 레코드를 설정해 `example.com`으로 오는 메일이 `mailhost.example.com.`로 모두 오도록 함
- 6행
  - `mydestination`에 `$mydomain`(example.com)을 설정
  - `user@example.com` 형식의 메일을 로컬에서 최종 배달
  - 일반 서버와 달리 `example.com` 도메인에 대한 메일을 수신

### mailbox 접근

- 사용자는 mailbox에 다음 방법으로 접근 가능
  - `/var/mail`에서 mailbox 디렉터리에 직접 접근
  - POP, IMAP 프로토콜을 사용해 접근
  - aliases 설정을 통해 특정 컴퓨터에 mailbox를 설정

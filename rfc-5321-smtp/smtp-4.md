# 4. The SMTP Specifications

## 4.1. SMTP Commands

### 4.1.1. Command Semantics and Syntax

- SMTP 명령어는 메일 전송, 메일 시스템 기능을 정의
- 메일 트랜잭션에는 명령의 매개 변수로 전달되는 여러 데이터가 존재
  - 역방향 경로 `reverse-path`: `MAIL` 명령어의 인수
  - 정방향 경로 `forward-path`: `RCPT` 명령어의 인수
  - 메일 데이터: `DATA`의 인수
- 메일 트랜잭션의 매개 변수들은 트랜잭션 완료 표시가 있기 전까지 유지되어야 하기 때문에 별도의 버퍼가 필요
  - `reverse-path` 버퍼, `forward-path` 버퍼, `mail data` 버퍼 존재
- 일부 명령어(`RSET`, `DATA`, `QUIT`)는 매개 변수를 허용하지 않음

#### 4.1.1.1 Extended HELLO (EHLO) or HELLO (HELO)

- SMTP 클라이언트를 SMTP 서버에 식별, 매개 변수로는 SMTP 클라이언트의 FQDN 사용
- SMTP 클라이언트
  - SMTP 클라이언트는 메일 트랜잭션 전 반드시 `HELO`나 `EHLO` 사용
- SMTP 서버
  - 해당 명령에 대해 연결 환영 응답 및 SMTP 클라이언트에게 자신을 식별
  - SMTP 서버는 `HELO` 명령을 반드시 지원
- 해당 명령과 응답 `250 OK`는 SMTP 클라이언트, 서버 모두 초기 상태이며 진행중인 트랜잭션이 없고 모든 버퍼와 상태 테이블이 비워졌음을 확인
- EHLO에 대한 응답은 지원되는 모든 확장에 대한 키워드를 반드시 포함해 여러 줄로 표시

#### 4.1.1.2. MAIL (MAIL)

- 메일 트랜잭션을 시작하는 데 사용
- 매개변수로 역방향 경로 `reverse-path`가 포함

```plain text
mail = "MAIL FROM:" Reverse-path [SP Mail-parameters] CRLF
```

#### 4.1.1.3. RECIPIENT (RCPT)

- 메일의 개별 수신자를 식별하는데 사용
  - 여러 수신자는 이 명령어를 여러 번 사용해 지정
  - 해당 명령어로 매개변수 정방향 경로를 `forward-path` 버퍼에 포함, 나머지 버퍼를 수정하지 않음
- 매개변수로 정방향 경로 `forward-path`가 포함

```plain text
rcpt = "RCPT TO:" ( "<Postmaster@" Domain ">" / "<Postmaster>" / Forward-path ) [SP Rcpt-parameters] CRLF
```

#### 4.1.1.4. DATA (DATA)

- SMTP 서버는 일반적으로 `DATA`에 `354` 응답을 보낸 후 `<CRLF>` 다음 줄부터 SMTP 클라이언트의 메일 데이터로 처리
- 서버는 메일 데이터 종료 표시 `<CRLF>.<CRLF>`를 수신하면 버퍼의 정보를 사용해 저장된 메일 트랜잭션 정보를 처리해야 함
  - 명령 완료 시 버퍼 정보 삭제
  - 처리 성공 시 SMTP 서버는 `OK` 응답, 실패 시 실패 응답을 보내야 함
  - 부분적인 성공과 실패는 허용되지 않음
  - 만약 SMTP 서버가 `250 OK`를 보냈다면 해당 메시지의 전송에 대해 완전한 책임을 지게 됨
- SMTP 서버가 메시지 릴레이 및 최종 전송 시 메시지 내용의 맨 위에 추적 기록(타임 스탬프, `Received` 라인)을 삽입
  - 해당 기록은 메시지를 보낸 호스트의 ID, 메시지를 수신한 호스트 ID, 메시지 수신 날짜 및 시간을 가짐
  - 릴레이된 메시지는 여러 개의 추적 기록을 가짐

```plain text
data = "DATA" CRLF
```

#### 4.1.1.5 RESET (RSET)

- 메일 트랜잭션 중단
  - 저장된 역방향 경로, 정방향 경로, 메일 데이터는 모두 폐기
  - 모든 버퍼 및 상태 데이터는 지워져야 함
- SMTP 서버는 인수 없는 RSET 명령에 대해 `250 OK` 응답을 반드시 보내야 함
- SMTP 클라이언트는 언제든지 이 명령을 사용 가능
- SMTP 서버는 `RSET` 수신 결과로 연결을 닫아서는 안됨

```plain text
rset = "RSET" CRLF
```

#### 4.1.1.6. VERIFY (VRFY)

- SMTP 서버에게 해당 명령의 매개변수가 사용자 또는 메일함을 식별하는지 확인 요청
- 해당 명령은 `reverse-path` 버퍼, `forward-path` 버퍼, `mail data` 버퍼에 영향을 미치지 않음

```plain text
vrfy = "VRFY" SP String CRLF
```

#### 4.1.1.7. EXPAND (EXPN)

- SMTP 서버에게 해당 명령의 매개변수가 메일링 리스트를 식별하는지 확인하고 그렇다면 해당 리스트의 멤버를 반환하도록 요청함
- 해당 명령은 `reverse-path` 버퍼, `forward-path` 버퍼, `mail data` 버퍼에 영향을 미치지 않음
- 해당 명령은 언제든지 사용 가능

```plain text
expn = "EXPN" SP String CRLF
```

#### 4.1.1.8. HELP (HELP)

- SMTP 서버가 SMTP 클라이언트에게 도움이 되는 정보를 보내도록 함
- 해당 명령은 `reverse-path` 버퍼, `forward-path` 버퍼, `mail data` 버퍼에 영향을 미치지 않음
- 해당 명령은 언제든지 사용 가능

```plain text
help = "HELP" [ SP String ] CRLF
```

#### 4.1.1.9. NOOP (NOOP)

- SMTP 서버는 해당 명령에 대해 `250 OK` 응답을 보내는 것 이외에 어떤 동작도 하지 않음
- 해당 명령은 `reverse-path` 버퍼, `forward-path` 버퍼, `mail data` 버퍼에 영향을 미치지 않음
- 해당 명령은 언제든지 사용 가능

```plain text
noop = "NOOP" [ SP String ] CRLF
```

#### 4.1.1.10. QUIT (QUIT)

- SMTP 서버는 해당 명령에 대해 `221 OK`를 보내고 전송 채널을 닫아야 함
  - 해당 명령을 수신하고 응답하기 전까지는 전송 채널을 닫아서는 안됨
- SMTP 클라이언트는 `QUIT` 명령을 보내기 전까지 전송 채널을 닫아서는 안되며 해당 명령에 대한 응답을 받을 때까지 기다려야 함
- 위 경우를 위반하거나 장애로 인해 연결이 조기에 닫힌 경우
  - 보류중인 트랜잭션은 모두 중단하고, 완료된 트랜잭션은 되돌려서는 안됨
  - 진행 중인 명령, 트랜잭션은 임시 오류를 수신한 것처럼 동작(`4XX`)
- 해당 명령은 언제든지 사용 가능
  - 명령 사용 시 완료되지 않은 모든 메일 트랜잭션이 중단

```plain text
quit = "QUIT" CRLF
```

### 4.2.1. Reply Code Severities and Theory

- 응답 코드의 세 자리는 각각 특별한 의미를 가짐
- 더 상세한 상태 코드 정보는 RFC 3463을 확인

#### 첫 번째 숫자

- 응답의 심각도를 표현하며, 클라이언트는 해당 숫자를 검토해 다음 동작을 결정할 수 있음
- 2yz: 긍정 완료
  - 명령이 수락되어 요청한 작업이 성공적으로 완료, 새로운 요청 가능
- 3yz: 긍정 중간
  - 명령이 수락되었지만 요청된 작업은 추가 정보 수신을 기다리며 보류 중
  - SMTP 클라이언트는 다른 명령을 사용해 추가 정보를 보내야 함
- 4yz: 일시적 부정 완료
  - 명령이 수락되지 않았고 요청된 작업이 발생하지 않음
  - 오류 조건이 일시적이므로 작업을 다시 요청할 수 있음
  - 동일한 명령 형식으로 서버와 클라이언트의 근본적 속성(주소 등)을 변경하지 않고 다시 시도했을 때 성공하면 4yz 응답임
  - SMTP 클라이언트는 동일한 요청을 다시 시도해야 함
- 5yz: 영구적 부정 완료
  - 명령이 수락되지 않았고 요청된 작업이 발생하지 않음
  - SMTP 클라이언트는 동일한 요청을 반복해서는 안됨

#### 두 번째 숫자

- x0z 구문
  - 구문 오류
  - 정의되지 않았거나 구현되지 않은 명령
  - 불필요한 명령
- x1z 정보
  - 상태 또는 도움말과 같은 정보 요청에 대한 응답
- x2z 연결2
  - 전송 채널을 나타내는 응답
- x3z, x4z 지정되지 않음
- x5z 메일 시스템
  - 수신 메일 시스템의 상태 표현

#### 세 번째 숫자

- 두 번째 숫자로 지정된 각 범주에서 더 세분화된 의미를 부여

#### 응답 텍스트

- 응답 텍스트는 권장 사항이지 필수 사항이 아니며 명령에 따라 변경 가능
- 여러 줄 응답에서 각 줄의 응답 코드는 모두 동일해야 함

```plain text
250-First line
250-Second line
250-234 Text beginning with numbers
250 The last line
```

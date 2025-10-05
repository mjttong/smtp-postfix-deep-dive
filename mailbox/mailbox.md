# mailbox

## postfix와 mailbox

- postfix는 별도의 MDA나 설정이 없다면 수신한 메일을 로컬 메일박스(`/var/mail`)에 저장
- 기본적으로 메일박스는 `mbox` 방식을 사용

### mbox, maildir

- mbox 방식
  - 사용자마다 자신의 이름으로 되어있는 파일이 하나씩 존재
  - 해당 파일에 자신에게 온 메일이 모두 저장됨
  - 모든 메일이 하나의 파일에 이어서 기록됨
- maildir 방식
  - 사용자마다 자신의 홈 디렉토리에 maildir 디렉토리 존재
  - 해당 디렉토리로 자신에게 온 메일 파일이 저장됨
  - 각각의 메일이 별도의 파일로 저장되고, 해당 파일이 디렉토리 구조를 통해 관리됨

### mailbox 관련 설정

- 메일박스 위치, 스타일
  - `home_mailbox`: `local` 사용자 홈 디렉토리 내의 메일박스 경로 지정
  - `mail_spool_directory`: `local` UNIX 스타일 메일박스가 보관되는 디렉토리 지정
- 메일박스 전송 명령 및 전송 방식
  - `mailbox_command`: `local` 전송 에이전트가 메일박스 전송에 사용할 외부 명령을 지정(전역적)
    - 유사한 파라미터 `mailbox_command_maps`(지역적, 수신자별 선택 가능)
  - `mailbox_transport`: `local` 전송 에이전트가 사용할 메시지 전송 방식을 지정(전역적)
    - 유사한 파라미터 `mailbox_transport_maps`(지역적, 수신자별 선택 가능)
- `mailbox_size_limit`: 메일박스 크기 제한
- 이 외에도 메일박스 잠금, 가상 메일박스 관련 파라미터 존재
- `mailbox` 관련 설정은 `local` 전달 에이전트와 관련
  - `pipe`는 메일박스를 관리하기보다는 메일을 외부 프로그램에 파이프하는 역할이기 때문에 무관

## 메일 전달 에이전트

### 전달 에이전트 간 관계

- Postfix는 메일 전달 시, 모든 전달 에이전트를 다 사용 가능하며 특정 전달 에이전트만 사용하는 것도 가능
  - 큐 관리자 `qmgr`가 다양한 전달 에이전트에게 전달 요청이 가능하기 때문
  - 파라미터 설정을 통해 유연한 라우팅이 가능하기 때문
- `local`이 직접, 간접적으로 `pipe`를 호출 가능
  - local과 pipe 기능을 연계해 사용하는 경우가 많음

### local

- 메일을 로컬 수신자(서버 내의 사용자 계정, 파일, 외부 명령)에게 전달
- 별칭 기반 라우팅, 고급 라우팅, 다른 전달 에이전트와의 연결 등 복잡한 메일 정책 구성에 유용

### pipe

- 메일을 외부 명령(외부 스크립트, 프로그램, DB, 메시지 브로커 등)으로 전달
- postfix와 외부 서비스를 연결하는 전달 에이전트

## pipe 설정 및 GCP queue와 연결

### lookup table

- postfix에서 액세스 제어, 컨텐츠 필터링, 라우팅 등에 조회 테이블(lookup table)을 사용
- 조회 테이블은 항상 type:table의 구조(예시: `hash:/etc/postfix/map`)
  - type에는 `hash`, `regexp` 등의 단순한 구조부터 `mysql`, `mongodb` 등의 외부 데이터베이스까지 지원
  - table에는 해당 테이블의 경로 작성
- `hash` 테이블을 사용해 메일 라우팅을 설정
  - 시스템 계정(root, daemon, postmaster)은 local이 배달
  - 그 외의 모든 계정으로부터 오는 메일은 pipe가 배달

### `transport_maps` 설정

```plain text
# /etc/postfix/transport
root    local:
daemon  local:
postmaster local:
```

```plain text
# /etc/postfix/main.cf
local_transport = gcppipe:
transport_maps = hash:/etc/postfix/transport
```

- `transport_maps`는 수신자 주소를 기반으로 라우팅을 결정
- 메일 송수신 시나리오
  - 발신 시 `gcppipe`로 메일이 전달되어서는 안됨
  - 수신 시 시스템 계정을 제외한 모든 계정의 메일이 `gcppipe`로 전달, `gcppipe`가 스크립트를 실행해 메일을 pub/sub으로 전송
- 메일 송신
  - `transport_maps`을 검사
  - 시스템 계정(root, daemon, postmaster) 메일의 경우 local로 전달
  - `transport_maps`에 명시되지 않은 계정은 smtp를 사용해 전송
- 메일 수신
  - `transport_maps`을 검사
  - 시스템 계정(root, daemon, postmaster) 메일의 경우 local로 전달
  - `transport_maps`에 명시되지 않은 모든 수신 메일은 `local_transport`의 기본값인
    `gcppipe`로 전달
  - `gcppipe`의 스크립트 실행

### pipe 설정

- type: 프로세스 간 통신(IPC) 방법 정의
  - inet(TCP/IP), unix(UNIX 소켓) 등의 값을 사용 가능
- private: 서비스가 내부에 위치하는지 여부
- unprivileged(unpriv): 서비스가 root(또는 postfix 소유자) 권한으로 실행되는지 여부
- chroot: 서비스가 chroot 환경에서 실행되는지
- wakeup: 서비스가 정기적으로 깨어나야 하는 시간 설정
- maxproc: 동시에 실행될 수 있는 최대 서비스 프로세스 수
- command + args: 실행될 프로그램의 경로, 명령줄 인수

```plain text
# /etc/postfix/master.cf
# service type  private unpriv  chroot  wakeup  maxproc command + args
gcppipe   unix  -       n       n       -       -       pipe 
  flags=FqX user=nobody argv=/usr/local/bin/gcppipe_script.go ${sender} ${recipient}
```

- pipe 전달 에이전트를 통해 메일을 `gcppipe_script.go` 스크립트로 전달

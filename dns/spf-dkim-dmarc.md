# SMTP 보안 레코드(SPF, DKIM, DMARC)

- SPF, DKIM는 수신 MTA가 메일을 받았을 때, 해당 메일이 스팸 서버로부터 사칭되거나 위조, 변조되지 않았는지 검증함
- DMARC는 SPF, DKIM 레코드를 확인한 후 수행할 작업을 수신 MTA에게 알려줌

## SPF(Sender Policy Framework)

- 특정 도메인 이름으로 메일을 보낼 수 있는 서버 IP 목록
  - 해당 서버는 특정 도메인을 SMTP 명령 `MAIL FROM`에서 사용 가능
- SPF 레코드를 통해 해당 IP 주소만 해당 도메인 이름으로 메일을 보낼 권한이 있다 선언
- 수신 서버는 SPF 레코드의 서버 목록을 확인해 발신 서버가 해당 도메인에 대해 권한이 있는지 확인

### SPF 레코드 구조

- 레코드 이름은 해당 도메인(`example.com`), 컨텐츠는 아래와 같은 값을 가짐

```plain text
v=spf1 ip4:1.1.1.1 ip4:2.2.2.2 include:test.com -all
```

- TXT 레코드의 이름으로 지정된 도메인 `example.com`에 대해 SPF 적용
- `v=spf1`: SPF 레코드라고 정의
- `ip4:1.1.1.1`: `example.com` 이름으로 메일 발신이 가능한 IP
- `include:test.com`: `test.com` 도메인의 SPF 레코드도 함께 참고
  - `example.com`과 `test.com` 양쪽에 허가된 모든 IP가 메일 발송 가능
- `-all`: SPF 레코드에 열거되지 않은 IP는 모두 차단
  - `~all`: SPF 레코드에 열거되지 않은 IP는 일단 허용하지만 스팸 분류
  - `+all`: SPF 레코드에 열거되지 않은 IP도 모두 허용(= 사실상 SPF 검사 안함)

### 사용 이유

- 도메인 사칭 방지
  - 공격자가 특정 도메인을 사칭해 스팸 발송 시 도메인 평판 하락
  - SPF 설정으로 인증된 IP만 메일 발신 가능하도록 제한해 도메인 평판 보호
- 메일 배달률 향상
  - Gmail, Outlook 등 주요 메일 서비스는 SPF 검증을 필수로 수행
  - SPF 미설정 도메인에서 발송된 메일은 스팸 분류 또는 반송 가능성 높음

## DKIM(DomainKeys Indentified Mail)

- 메일 내용의 변조 여부를 검증
- 발신 서버는 공개키와 비밀키를 생성
  - 공개키는 발신 서버 도메인의 TXT 레코드에 저장(DKIM 레코드)
  - 비밀키는 발신 서버만 보관
- 발신 서버가 메일 전송 시 디지털 서명과 메일 해시 값을 메일 헤더에 포함해 함께 전송
- 수신 서버는 메일을 받아 발신 서버의 DKIM 레코드 조회, 공개 키 획득
  - 공개 키로 디지털 서명 검증
  - 메일 해시 값을 통해 메일 변조 여부 검증

### DKIM 레코드 구조

```plain text
[selector]._domainkey.[domain]
```

- DKIM 레코드는 특수한 이름으로 저장
- `[selector]`: dkim 레코드를 식별하는 값, 사용자 지정 값이 됨
- `._domainkey.`는 모든 DKIM 레코드 이름에 포함됨
- `[domain]`: DKIM 레코드를 적용할 도메인

```plain text
v=DKIM1; p=[공개키 값]
```

- DKIM 컨텐츠에서 공개 키를 기록
- `v=DKIM1`: DKIM 레코드라고 정의
- `p=`: DKIM 공개 키 값을 기록

### DKIM 헤더

```plain text
v=1; a=rsa-sha256; 
d=example.com; s=big-email;
h=from:to:subject;
bh=uMixy0BsCqhbru4fqPZQdeZY5Pq865sNAnOAxNgUS0s=;
b=LiIvJeRyqMo0gngiCygwpiKphJjYezb5kXBKCNj8DqRVcCk7obK6OUg4o+EufEbB
tRYQfQhgIkx5m70IqA6dP+DBZUcsJyS9C+vm2xRK7qyHi2hUFpYS5pkeiNVoQk/Wk4w
ZG4tu/g+OA49mS7VX+64FXr79MPwOMRRmJ3lNwJU=
```

- `v=`: DKIM 버전
- `a=`: 디지털 서명과 해시에 사용된 알고리즘
- `d=`: 발신 도메인
- `s=`: 수신 MTA가 사용해야 할 selector
- `h=`: 서명에 포함된 헤더 필드 목록
- `bh=`: 메일 본문의 해시값 (본문 변조 검증용)
- `b=`: 디지털 서명값 (개인키로 생성)

### DKIM 검증 시나리오

- 메일 수신 서버는 메일 헤더에서 `d`(domain)와 `s`(selector)를 추출
- DNS에서 `s._domainkey.d`을 통해 DKIM 레코드 조회, DKIM 공개 키 획득
- 공개키로 `b` 서명 검증 및 `bh` 본문 해시 확인

## DMARC(Domain-based Message Authentication, Reporting, and Conformance)

- 메일을 SPF, DKIM 레코드로 검증한 후 실패한 메일에 대해 수행할 작업(스팸 분류, 차단 등)을 결정
- 발신 서버가 자신의 도메인으로 온 메일에 대해 SPF, DKIM이 실패할 경우 수신 서버가 어떻게 행동해야 할 지를 정해주는 것
  - 수신 서버는 DMARC 정책을 따를 의무는 없음, 이는 권장 사항임

### DMARC 레코드 구조

```plain text
v=DMARC1; p=quarantine; adkim=s; aspf=s;
```

- `v=DMARC1`: DMARC 레코드 버전
- `p=`: SPF, DKIM 실패 시 수행할 작업
  - `none`: 통과 허용
  - `quarantine`: 스팸함 전송(격리)
  - `reject`: 차단
- `adkim=`, `aspf=`: SPF, DKIM을 얼마나 엄격하게 검사할지 결정
  - `r`: relaxed, 완화
  - `s`: strict, 엄격

## 메일 보안 레코드 설정

- 보안 레코드 설정
  - MTA가 메일을 발신하는 경우, 신뢰성 보장 및 스팸 메일 분류를 막기 위해 설정해야 함
- 보안 레코드 관련 프로그램 설치
  - SPF: 메일 수신을 위해 필요
    - 메일 수신 시 발신 서버의 SPF 레코드를 검증하기 위해 필요
  - DKIM: 메일 발신, 수신을 위해 필요
    - 메일 발신 시 해당 메일에 DKIM 헤더를 추가해야 하기 때문에 필요
    - 메일 수신 시 DKIM 헤더를 바탕으로 DKIM 레코드를 조회한 후 서명을 검증하는 과정을 거쳐야 하기 때문에 필요
  - DMARC: 메일 수신을 위해 필요
    - 메일 수신 시 SPF, DKIM 검사를 거친 후 DMARC 레코드 조회를 통해 수행할 작업을 결정하기 위해 필요

### DNS 레코드 설정

- 발신 서버에서 SPF, DKIM, DMARC 레코드를 등록
  - SPF 레코드: 도메인으로 메일을 보낼 수 있는 IP 명시
  - DKIM 레코드: 메일 서명 검증용 공개키 게시
  - DMARC 레코드: SPF/DKIM 실패 시 처리 방침 명시

### Postfix 서버 프로그램 설정

- 발신 서버
  - openDKIM: 발신 메일에 DKIM 서명 추가
- 수신 서버
  - postfix-policyd-spf-python: 발신 서버 IP와 SPF 레코드 비교 검증
  - openDKIM: DKIM 서명, 해시 값 검증
  - openDMARC: SPF, DKIM 결과와 DMARC 정책을 바탕으로 처리

### mailgun과 메일 보안 레코드

- 발신의 경우
  - mailgun에서 제공하는 SPF, DKIM 레코드 추가, DMARC 레코드 설정
  - mailgun이 발신 릴레이로 역할하기 때문에 별도로 openDKIM을 설치하지 않아도 알아서 DKIM 처리 가능
- 수신의 경우
  - mailgun을 오직 발신 릴레이 용도로만 쓰기 때문에 mailgun과 무관
  - SPF, DKIM, DMARC 프로그램 설치 및 설정 필요

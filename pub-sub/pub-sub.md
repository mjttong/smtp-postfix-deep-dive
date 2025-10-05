# Pub/Sub

## Message Queue 사용

- 메일 컴포넌트(MUA, MTA, MDA)는 각각 독립적으로 교체 가능한 모듈 구조
  - MUA(Outlook, Thunderbird), MTA(Postfix, Qmail) 등을 자유롭게 조합 가능
- 해당 프로젝트에서는 메일 컴포넌트의 낮은 결합도, 높은 유연성을 구현하기 위해 메시지 큐 사용

### 목적

- MTA-MDA 간 결합도 최소화
  - MTA는 메시지 큐에 메일 전달만 수행
  - MDA는 큐에서 메일을 가져와 독립적으로 처리
  - 서로의 구현 방식에 영향을 받지 않으며, '메시지 큐로 메일 전달'이라는 로직으로 쉽고 간편하게 개발 가능
- 메시지 큐를 통한 메일 유실 방지
- 각 컴포넌트가 자신의 처리 속도에 맞춰 메일 소비 가능

## Pub/Sub vs Cloud Tasks

🔗 [Pub/Sub 또는 Cloud Tasks 선택](https://cloud.google.com/pubsub/docs/choosing-pubsub-or-cloud-tasks?hl=ko)

- 원하는 큐 서비스의 특성
  - MTA는 메시지 큐에 메일을 전달하는 행위만 수행
  - MDA가 메시지 소비 방식을 자유롭게 선택
  - MTA, MDA가 서로의 존재나 구현 방식을 알 필요 없음, 서로의 행위에 간섭하지 않음
  - AWS SQS와 유사한 특성을 갖는 큐 서비스

### Cloud Tasks 특성

- 게시자가 명시적으로 소비자를 호출
  - 게시자가 실행에 대해 완전한 제어를 유지
  - 소비자의 메시지 소비 행동(실행 시점, 재시도 정책 등)을 직접적으로 제어
- 게시자가 특정 엔드포인트를 지정해 태스크 전달
- 따라서 각 컴포넌트 간 결합도 증가, 긴밀하게 연결되어야 해 해당 프로젝트에 사용하기 부적절

### Pub/Sub 특성

- 게시자는 토픽에 메시지만 게시, 구독자의 실행을 암시적으로 유지
- 구독자는 독립적으로 메시지 소비
- AWS SQS와도 유사한 동작 방식
- 따라서 각 컴포넌트는 서로의 정보가 불필요, 결합도 감소로 해당 프로젝트에 사용하기 매우 적절
- 해당 프로젝트에서는 메시지 큐로 Pub/Sub을 사용하기로 결정

## Postfix + Pub/Sub

### 발신 시나리오

MUA → Pub/Sub → go 스크립트 → postfix → mailgun → 최종 목적지

- MUA가 메일을 pubsub 토픽에 게시
- go 스크립트가 pubsub에서 메시지 pull
- go의 smtp 패키지로 Postfix에 SMTP 연결(587 submission)
- Postfix가 mailgun 릴레이를 통해 외부 발송

### 수신 시나리오

외부 MTA → postfix → pipe (gcppipe) → go 스크립트 → Pub/Sub → MDA

- 외부 MTA가 25번 포트로 메일 발송
- Postfix가 메일 수신
- 시스템 계정 메일은 local 에이전트로, 그 외의 모든 메일은 gcppipe 에이전트로 전달
- gcppipe 스크립트가 메일을 pubsub에 게시
- MDA가 pubsub에서 메일 소비 및 처리

### submission 수정

- go 스크립트를 통해 `localhost:587`로 SMTP 연결 시 설정해둔 Dovecot SASL 인증을 사용하지 않음
  - `mynetworks`에 localhost가 포함되어 있어 postfix가 localhost를 신뢰하기 때문에 인증하지 않음
- 로컬(go 스크립트 → postfix)에 Dovecot SASL인증을 추가하지 않고 그대로 설정 유지
  - go 스크립트와 postfix는 같은 호스트에서 실행, 인증을 추가하는 것이 더 비효율적
  - 외부 클라이언트는 여전히 postfix 접근 시 Dovecot SASL 인증 필요

# Submission

## MSA(Message Submission Agent)

### Submission과 Transfer의 분리

- 메시지 제출(submission)과 메시지 전송(transfer)을 분리
- MUA로부터 메시지 제출을 받는 프로세스를 MSA(Message Submission Agent)라고 정의
- MSA는 MTA에게 메시지를 제출, MTA는 메시지를 다른 MTA에게 전송

### 도입 배경

- 보안 중요성 증대
  - 초기 메시지 제출 시 인증 및 권한 부여의 중요성 증가
  - 많은 사이트들이 25번 포트의 아웃바운드 트래픽 금지
- 메시지 불완전성
  - 제출되는 메시지가 완성된 상태가 아닐 수 있음
  - 메시지를 표준에 맞게 보완할 필요가 있음
- 사이트 정책 준수
  - 사이트의 정책에 따라 메시지 텍스트를 검사하거나 수정해야할 수 있음
  - 이러한 메시지의 수정은 MTA 기능의 영역을 벗어나는 것으로 간주

### 587번 포트

- submission을 위해 특별히 예약된 포트
- SASL(Simple Authentication and Security Layer) 사용 의무, SSL/TLS 인증 권장

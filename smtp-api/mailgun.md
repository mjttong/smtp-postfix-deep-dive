# mailgun API

## 25번 포트 아웃바운드 차단

> 🔗 [GCP: 인스턴스에서 이메일 보내기](https://cloud.google.com/compute/docs/tutorials/sending-mail?hl=ko)
>
> 악용의 위험으로 인해 대상이 VPC 네트워크 외부에 있을 때 대상 TCP 포트 25로의 연결은 차단됩니다.

- 대부분의 VPS(Virtual Private Server)는 25번 포트의 아웃바운드 트래픽을 차단함
  - VM이 스팸 메일 서버로 악용될 수 있기 때문
  - 25번 포트 인바운드는 허용 가능
- 요청은 아웃바운드(패킷 전송), 응답은 인바운드(패킷 수신)
  - 요청 시 소스 포트는 랜덤 값 포트, 목적지 포트 특정 값(특정한 서비스) 포트
  - 즉, 25번 포트의 아웃바운드 트래픽을 차단 = 목적지 포트가 25번인 패킷을 모두 차단
  - 내 VM(랜덤 포트)→ 외부 서버(25번 포트) 방향으로 가는 패킷은 네트워크 레벨에서 항상 드랍
- 그렇다면 25번 포트가 아닌 다른 포트를 목적지 포트로 해서 메일을 전송한다면?
  - 외부 메일 서버의 방화벽에 의해 패킷 차단
  - 외부 MTA는 표준적으로 25번 포트에서만 SMTP 연결을 수신, 다른 포트로 전송 시 외부 메일 서버가 수신 대기 중이지 않아 연결 실패
  - 587번(submission)은 클라이언트 → MSA 전용 포트로 MTA 간 통신에는 사용되지 않음

### 해결 방법

- 25번 포트가 뚫려있는 VPS를 사용
  - GCP와 복잡한 연동
  - 제공 업체를 찾기 힘듦
  - 추가적인 비용이 듦
- SMTP API를 SMTP 릴레이 서버로 사용
  - Postfix(아웃바운드) → SMTP API(릴레이 용도) → 외부 MTA
  - GCP에서 즉시 적용해 사용 가능
  - SMTP API는 발신 릴레이로만 사용
  - SMTP API로는 Mailgun 사용

### SMTP API 사용

- SMTP API를 사용한다면 해당 API가 Postfix를 대체할 수 있는데, Postfix를 써야 하는 이유가 있는지?
- 시스템 결합도 감소
  - MUA는 오직 Postfix와만 통신하고, Postfix가 외부 SMTP API와 통신
  - 애플리케이션과 외부 API 간의 결합도가 매우 낮아짐
- 확장성 및 유연한 라우팅
  - 특정 SMTP API에 종속되지 않는 아키텍처이므로 확장성이 높음
  - Postfix의 강력한 라우팅 기능을 활용해 상황에 따라 다른 메일 전송 경로를 설정할 수 있음
- 안정적인 큐 관리 및 로깅
  - Postfix는 메일 전송에 실패하더라도 메일을 큐에 보관하고 재전송을 시도
  - 이 덕분에 외부 API나 네트워크에 일시적인 문제가 생겨도 메일이 유실되지 않고 시스템의 안정성을 보장
  - 모든 전송 내역이 Postfix 로그에 기록되어 중앙 집중적인 디버깅과 감사 기록을 제공

## mailgun 사용하기

- mailgun, postfix, dns 서버 설정 필요
- 설정 후 메일을 전송 시, postfix에서 mailgun을 릴레이 서버로 사용해 메일이 발신됨(발신 릴레이)
- MX 레코드는 mailgun 서버를 경유하지 않고 Posifix로 직접 도착하도록 설정
  - mailgun을 발신 릴레이 용도로만 사용하기 위함
  - Postfix → pipe → 스크립트 → Pub/Sub로 전달

### Postfix 설정

🔗 [Mailgun: SMTP Relay](https://documentation.mailgun.com/docs/mailgun/user-manual/smtp-protocol/smtp-relay)

```plain text
# /etc/postfix/main.cf:

relayhost = [smtp.mailgun.org]:587
smtp_sasl_auth_enable = yes
smtp_sasl_password_maps = hash:/etc/postfix/sasl_passwd
smtp_sasl_security_options = noanonymous

smtp_tls_note_starttls_offer = yes
```

```plain text
# /etc/postfix/sasl_passwd:
[smtp.mailgun.org]:587 your_mailgun_smtp_user@your_subdomain_for_mailgun:your_mailgun_smtp_password
```

- Postfix(MTA)가 Mailgun(발신 릴레이)에 연결하기 위해 587번 포트 사용
- 587번 포트 사용을 위해 SASL 설정 필수, TLS 설정 권장

### mailgun의 역할

- SMTP 릴레이로써 장점
  - 메일 로그, 지표 확인이 편리
  - 자동 SPF/DKIM 설정으로 보안 구성 간소화
  - 도메인 평판 관리 및 배달 성공률 향상
- 프로젝트에서의 제한적 사용
  - Mailgun은 완전한 메일 솔루션을 제공하지만, 이 프로젝트에서는 SMTP 릴레이 기능만 활용
  - Postfix MTA 구축이라는 프로젝트 주제를 유지하면서도 25번 포트 제약을 우회하는 해결책으로만 사용하고자 함

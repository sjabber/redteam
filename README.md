# redteam - 악성메일 대응 솔루션 (핵심 코드 제외됨)
- RedTeam 사이트 주소 : http://mert.koreacentral.cloudapp.azure.com:8080/

## 프로젝트 개요

#### 프로젝트 소개

보안기술의 발달로 시스템의 보안성은 강화되고 있는 반면, 상대적으로 사람들의 보안의식은 취약한 것으로 나타난다.<br>
이런 사람의 취약점을 공략하여 원하는 정보를 얻는 사회공학적 해킹기법들의 증가로 보안 인식 제고를 위한<br>
무료 악성메일 대응 솔루션을 개발하여 소상공인 및 중소기업을 지원한다.

#### 시스템 구성도

<p align="center">
  <img src="https://user-images.githubusercontent.com/46473153/120894331-cfbfbf00-c652-11eb-9133-bb91904be68b.png">
</p>

<!-- ![image](https://user-images.githubusercontent.com/46473153/111221494-e7704400-861d-11eb-9db8-893311c6f144.png) -->

#### 시스템 개요도

![image](https://user-images.githubusercontent.com/46473153/111221598-0cfd4d80-861e-11eb-85dc-908e7276eaf4.png)

#### 시스템 주요 기능 리스트

![image](https://user-images.githubusercontent.com/46473153/111221652-1f778700-861e-11eb-9655-ed90bda56809.png)

#### 프로젝트 개발환경

![image](https://user-images.githubusercontent.com/46473153/111221992-944ac100-861e-11eb-9ab5-8a0cbb5b17fc.png)

## 작품의 주요 기능 및 장점

#### 기능

- 관리자 회원가입, 로그인, 정보수정
- 악성메일 템플릿 관리
- 훈련대상자 등록 및 관리
- SMTP 연결테스트 및 설정
- 악성메일 훈련 프로젝트 생성 및 예약
- 훈련 진행상황 확인, 훈련결과 확인

#### 장점

- 무료 제공으로 인해 중소기업들의 부담 없는 사용이 가능
- 훈련 대상자 관리가 편리함 (일괄등록, 삭제, 태그를 통한 분류기능 제공)
- 훈련 유형을 모의훈련 담당자가 자유롭게 수정이 가능함.
- 훈련 진행상황이나 결과를 확인이 가능하며 재훈련이 필요한 대상자들을 태그로 재등록하여 관리할 수 있음.

## 업무분장

#### 프로젝트 매니저

- 홍대화 수석님

#### 프론트엔드 담당 && 프로젝트 리더

- 김민석 (node.js - javascript)

#### 백엔드 담당 && 멘티

- 김태호 (gin framework - Go)
- 김광채 (Spring - java)

## Spring-framework 사용법

- springboot 디렉토리 이동 후 <br>
  java -jar mer-0.0.1-SNAPSHOT.jar 타이핑

- 이전 명령어 <br>
  java spring-framework 실행 방법 <br>
  java -Xmx100m -jar spring.jar

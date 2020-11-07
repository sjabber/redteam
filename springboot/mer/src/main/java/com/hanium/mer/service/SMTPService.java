package com.hanium.mer.service;

import com.hanium.mer.repogitory.SMTPRepository;
import com.hanium.mer.vo.SmtpVo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import javax.mail.*;
import java.time.LocalDateTime;
import java.util.Optional;
import java.util.Properties;

@Service
public class SMTPService {

    @Autowired
    SMTPRepository smtpRepository;

    public Optional<SmtpVo> getSMTP(Long user_no){
        Optional<SmtpVo> smtp =  smtpRepository.findByUserNo(user_no);
        return smtp;
    }

    public void setSMTP(Long user_no, SmtpVo newSmtp){
        SmtpVo smtp =  smtpRepository.findByUserNo(user_no).get();
        smtp.setSmtpHost(newSmtp.getSmtpHost());
        smtp.setSmtpPort(newSmtp.getSmtpPort());
        smtp.setSmtpProtocol(newSmtp.getSmtpProtocol());
        smtp.setSmtpTls(newSmtp.getSmtpTls());
        smtp.setSmtpTimeOut(newSmtp.getSmtpTimeOut());
        smtp.setSmtpId(newSmtp.getSmtpId());
        smtp.setSmtpPw(newSmtp.getSmtpPw());
        smtp.setModify(LocalDateTime.now());
        System.out.println(smtpRepository.save(smtp));
    }

    public void connectCheck(SmtpVo smtp_info) throws NoSuchProviderException, MessagingException {
        String host = smtp_info.getSmtpHost(); // 네이버일 경우 네이버 계정, gmail경우 gmail 계정
        String user = smtp_info.getSmtpId()+"@naver.com"; // 패스워드
        String password = smtp_info.getSmtpPw();

        // SMTP 서버 정보를 설정한다.
        Properties props = new Properties();
        props.put("mail.smtp.user", user);
        props.put("mail.smtp.host", host);
        props.put("mail.smtp.port", smtp_info.getSmtpPort());
        props.put("mail.smtp.startttls.enable", smtp_info.getSmtpTls());
        props.put("mail.smtp.auth", "true");
        props.put("mail.smtp.socketFactory.port",smtp_info.getSmtpPort());
        props.put("mail.smtp.socketFactory.class", "javax.net.ssl.SSLSocketFactory");
        props.put("mail.smtp.socketFactory.fallback", "false");

        Session session = Session.getDefaultInstance(props, new javax.mail.Authenticator() {
            protected PasswordAuthentication getPasswordAuthentication() {
                return new PasswordAuthentication(user, password);
            }
        });

        Transport tr = session.getTransport("smtp");
        tr.connect();

        /* 테스트 완료
        //메일보내는 부분 sendMessage(List InternetAddress)
        //InternetAddress를 배열로 받아서 후에 프로젝트에서 target의 메일들을 받아보내면 됨
        // setText 대신 템플릿(첨부파일, html 태그 포함된 글)을 보낼 수 있음
        try {
            MimeMessage message = new MimeMessage(session);
            message.setFrom(new InternetAddress(user));
            //수신자
            message.addRecipient(Message.RecipientType.TO, new InternetAddress("kimkc5215@naver.com"));

            // 메일 제목
            message.setSubject("SMTP TEST1111");

            // 메일 내용
            message.setText("Success!!");

            // send the message
            Transport.send(message);
            System.out.println("Success Message Send");

        } catch (Exception e) {
            System.out.println("error");
            e.printStackTrace();
            return false;
        }
        return true;
        */

    }

}

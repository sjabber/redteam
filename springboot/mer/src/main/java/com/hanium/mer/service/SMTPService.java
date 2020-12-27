package com.hanium.mer.service;

import com.hanium.mer.repogitory.SMTPRepository;
import com.hanium.mer.vo.SmtpVo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.mail.javamail.JavaMailSenderImpl;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.Optional;
import java.util.Properties;

@Service
public class SMTPService {

    @Autowired
    SMTPRepository smtpRepository;

    @Autowired
    AESService aesService;

    public Optional<SmtpVo> getSMTP(Long user_no){
        Optional<SmtpVo> smtp =  smtpRepository.findByUserNo(user_no);
        try {
            //TODO LOG로 변경
            System.out.println(aesService.decAES(smtp.get().getSmtpPw()));//log
        }catch(Exception e){
            e.printStackTrace();
        }
        return smtp;
    }

    public void setSMTP(Long user_no, SmtpVo newSmtp) throws Exception{
        SmtpVo smtp =  smtpRepository.findByUserNo(user_no).get();
        smtp.setSmtpHost(newSmtp.getSmtpHost());
        smtp.setSmtpPort(newSmtp.getSmtpPort());
        smtp.setSmtpProtocol(newSmtp.getSmtpProtocol());
        smtp.setSmtpTls(newSmtp.getSmtpTls());
        smtp.setSmtpTimeOut(newSmtp.getSmtpTimeOut());
        smtp.setSmtpId(newSmtp.getSmtpId());
        smtp.setSmtpPw(aesService.encAES(newSmtp.getSmtpPw()));//
        System.out.println(smtp.getSmtpPw());//
        smtp.setModify(LocalDateTime.now());
        System.out.println(smtpRepository.save(smtp));
    }

    public void connectCheck(SmtpVo smtp_info) throws Exception {
        //TODO id(email 아이디) 유효성 검사하기
        JavaMailSenderImpl mailSender = new JavaMailSenderImpl();
        mailSender.setHost(smtp_info.getSmtpHost());
        mailSender.setPort(Integer.parseInt(smtp_info.getSmtpPort()));

        mailSender.setUsername(smtp_info.getSmtpId());
        mailSender.setPassword(aesService.decAES(smtp_info.getSmtpPw()));

        Properties props = mailSender.getJavaMailProperties();
        props.put("mail.smtp.startttls.enable", smtp_info.getSmtpTls().equals("1")?"true":"false");//smtp_info.getSmtpTls().equals("1")?"true":"false"
        props.put("mail.smtp.auth", "true"); // 인증이 필요없다면 불필요
        props.put("mail.transport.protocol", "smtp");
        props.put("mail.smtp.connectiontimeout", smtp_info.getSmtpTimeOut());
        // 1:"smtp" 2: "smtps"
        if( smtp_info.getSmtpProtocol().equals("2") || smtp_info.getSmtpPort().equals("465")){
            //SSL을 사용하지 않으면 필요없는 것 같다. SSL 사용포트 465, TLS는 587
            props.put("mail.smtp.socketFactory.port", smtp_info.getSmtpPort());
            props.put("mail.smtp.socketFactory.class", "javax.net.ssl.SSLSocketFactory");
            props.put("mail.smtp.socketFactory.fallback", "false");
        }


        mailSender.testConnection();

        /*
        //메세지 보내기: html 형태(검색시 금방 나옴), 파일도 보낼 수 있다.
        //여러 사람에게 보내기 List InternetAddress 형태? 검색필요
        //InternetAddress를 배열로 받아서 후에 프로젝트에서 target의 메일들을 받아보내면 됨
        MimeMessage message = mailSender.createMimeMessage();
        message.setFrom(new InternetAddress(smtp_info.getSmtpId()));
        //수신자
        message.addRecipient(Message.RecipientType.TO, new InternetAddress("kimkc5215@naver.com"));

        // 메일 제목
        message.setSubject("SMTP TEST1111");

        // 메일 내용
        message.setText("Success!!");
        mailSender.send(message);
        */

    }

}

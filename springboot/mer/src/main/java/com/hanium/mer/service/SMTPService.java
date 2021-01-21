package com.hanium.mer.service;

import com.hanium.mer.repogitory.SMTPRepository;
import com.hanium.mer.vo.SmtpVo;
import com.hanium.mer.vo.TargetVo;
import com.hanium.mer.vo.TemplateVO;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.mail.javamail.JavaMailSenderImpl;
import org.springframework.stereotype.Service;

import javax.mail.Message;
import javax.mail.internet.InternetAddress;
import javax.mail.internet.MimeMessage;
import java.time.LocalDateTime;
import java.util.Optional;
import java.util.Properties;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

@Slf4j
@Service
public class SMTPService {

    @Autowired
    SMTPRepository smtpRepository;

    @Autowired
    AESService aesService;

    public Optional<SmtpVo> getSMTP(Long user_no){
        Optional<SmtpVo> smtp =  smtpRepository.findByUserNo(user_no);
        try {
            log.info(aesService.decAES(smtp.get().getSmtpPw()));//log
        }catch(Exception e){
            e.printStackTrace();
        }
        return smtp;
    }

    public Boolean setSMTP(Long user_no, SmtpVo newSmtp) throws Exception{
        SmtpVo smtp =  smtpRepository.findByUserNo(user_no).get();
        smtp.setSmtpHost(newSmtp.getSmtpHost());
        smtp.setSmtpPort(newSmtp.getSmtpPort());
        smtp.setSmtpProtocol(newSmtp.getSmtpProtocol());
        smtp.setSmtpTls(newSmtp.getSmtpTls());
        smtp.setSmtpTimeOut(newSmtp.getSmtpTimeOut());
        String namePattern = "^[_a-z0-9-]+(.[_a-z0-9-]+)*@(?:\\w+\\.)+\\w+$";
        Matcher matcher = Pattern.compile(namePattern).matcher(newSmtp.getSmtpId());
        if(!matcher.matches()) {
            return false;
        }
        smtp.setSmtpId(newSmtp.getSmtpId());
        smtp.setSmtpPw(aesService.encAES(newSmtp.getSmtpPw()));
        log.debug(smtp.getSmtpPw());
        smtp.setModify(LocalDateTime.now());
        System.out.println(smtpRepository.save(smtp));
        return true;
    }

    //TODO 리턴값 변경 for id 유효성
    public Boolean connectCheck(SmtpVo smtp_info) throws Exception {
        JavaMailSenderImpl mailSender = new JavaMailSenderImpl();
        mailSender.setHost(smtp_info.getSmtpHost());
        mailSender.setPort(Integer.parseInt(smtp_info.getSmtpPort()));

        String namePattern = "^[_a-z0-9-]+(.[_a-z0-9-]+)*@(?:\\w+\\.)+\\w+$";
        Matcher matcher = Pattern.compile(namePattern).matcher(smtp_info.getSmtpId());
        if(!matcher.matches()) {
            return false;
        }
        mailSender.setUsername(smtp_info.getSmtpId());
        //파라미터 스트링 값으로 오니 decrypt 필요x
        mailSender.setPassword(smtp_info.getSmtpPw());

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
        return true;
    }

    public String sendMail(SmtpVo smtp_info, TargetVo target, TemplateVO template) throws Exception {
        //TODO checkConect와 중복제거
        JavaMailSenderImpl mailSender = new JavaMailSenderImpl();
        mailSender.setHost(smtp_info.getSmtpHost());
        mailSender.setPort(Integer.parseInt(smtp_info.getSmtpPort()));

        mailSender.setUsername(smtp_info.getSmtpId());
        //비번 초기 아이디 생성시, smtp 가 자동 생성되며 비밀번호에 BCrypt 문자가 생성
        // => Illegal base64 character 24 예외
        //TODO 예외 잡을지 말지
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

        log.info("smtp info {}", smtp_info);
        //메세지 보내기: html 형태(검색시 금방 나옴), 파일도 보낼 수 있다.
        //여러 사람에게 보내기 List InternetAddress 형태? 검색필요
        //InternetAddress를 배열로 받아서 후에 프로젝트에서 target의 메일들을 받아보내면 됨
        MimeMessage message = mailSender.createMimeMessage();
        //todo sender getId고치기+naver.com
        message.setFrom(new InternetAddress(smtp_info.getSmtpId()));// + "@" + smtp_info.getSmtpHost().substring(7)));
        //수신자
        log.info("target Email {}", target.getTargetEmail());
        message.addRecipient(Message.RecipientType.TO, new InternetAddress(target.getTargetEmail()));

        // 메일 제목
        message.setSubject(template.getMailTitle());

        // 메일 내용
        message.setContent(template.getMailContent(), "text/html; charset=utf-8");
        //message.setText(template.getMailContent());
        try {
            mailSender.send(message);
        }catch (Exception e){
            e.printStackTrace();//log.info("mail fail {}", e.printStackTrace());
            return "fail";
        }

        return "success";
    }

}

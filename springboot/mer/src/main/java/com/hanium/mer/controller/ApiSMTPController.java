package com.hanium.mer.controller;

import com.hanium.mer.TokenUtils;
import com.hanium.mer.service.AESService;
import com.hanium.mer.service.SMTPService;
import com.hanium.mer.vo.SmtpVo;
import com.hanium.mer.vo.Smtp_setting;
import io.jsonwebtoken.Claims;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import javax.mail.AuthenticationFailedException;
import javax.mail.MessagingException;
import javax.servlet.http.HttpServletRequest;
import java.io.UnsupportedEncodingException;
import java.util.Optional;

@RestController
public class ApiSMTPController {

    @Autowired
    SMTPService smtpService;
    @Autowired
    AESService aesService;

    @GetMapping("/setting/smtpSetting")
    public ResponseEntity<Object> getSMTPSetting(HttpServletRequest request) throws UnsupportedEncodingException {

        Optional<SmtpVo> smtp;

        Claims claims = TokenUtils.getClaimsFormToken(request.getCookies());
        if(claims != null){
            smtp = smtpService.getSMTP(Long.parseLong(claims.get("user_no").toString()));;
            return new ResponseEntity<Object>(new Smtp_setting(smtp.get()), HttpStatus.OK);
        }

        return new ResponseEntity<Object>("error", HttpStatus.FORBIDDEN);
    }

    //todo 비밀번호 암호화
    @PostMapping("/setting/smtpSetting")
    public ResponseEntity<Object> setSTMPSetting(HttpServletRequest request, @RequestBody SmtpVo newSmtp)
            throws UnsupportedEncodingException {

        Optional<SmtpVo> smtp;
        Claims claims = TokenUtils.getClaimsFormToken(request.getCookies());
        if (claims != null) {
            try{
                if ( !smtpService.setSMTP(Long.parseLong(claims.get("user_no").toString()), newSmtp)){
                    return new ResponseEntity<Object>("ID를 이메일 형식으로 입력해주세요.", HttpStatus.BAD_REQUEST);
                }
                System.out.println(newSmtp.toString());
                return new ResponseEntity<Object>(newSmtp.toString(), HttpStatus.OK);
            }catch(Exception e){
                e.printStackTrace();
                //에러 400-> smtp 정보확인
                //401 비밀번호확인 제대로 설정하기
                return new ResponseEntity<Object>("smtp 정보를 확인해주세요.", HttpStatus.BAD_REQUEST);
            }
        }

        return new ResponseEntity<Object>("토큰을 확인해주세요", HttpStatus.FORBIDDEN);
    }

    //todo 인증에러는 front 401은 connect에서, connect에서 auth 에러 시 아래 에러 정리하기
    @PostMapping("/setting/smtpConnectCheck")
    public ResponseEntity<Object> connectSTMPTest(@RequestBody SmtpVo smtp){
        try {
            if(!smtpService.connectCheck(smtp)){

                return new ResponseEntity<Object>("ID는 이메일 형식으로 입력해주세요.", HttpStatus.UNAUTHORIZED);
            }
        }catch (IllegalStateException e) {
            e.printStackTrace();
            return new ResponseEntity<Object>("이미 연결이 되어있습니다.", HttpStatus.BAD_REQUEST);
        } catch(AuthenticationFailedException e){
            e.printStackTrace();
            return new ResponseEntity<Object>("로그인 정보가 올바르지 않습니다.", HttpStatus.UNAUTHORIZED);
        }catch (MessagingException e){
            e.printStackTrace();
            return new ResponseEntity<Object>("SMTP 설정을 다시 확인해주세요.", HttpStatus.BAD_REQUEST);
        }catch(Exception e){
            e.printStackTrace();
            return new ResponseEntity<Object>("SMTP 설정을 다시 확인해주세요.", HttpStatus.BAD_REQUEST);
        }

        return new ResponseEntity<Object>("성공", HttpStatus.OK);
    }

    //project 생성 시
    @GetMapping("/api/smtpConnectSimpleCheck")
    public ResponseEntity<Object> connectionSTMPTest(HttpServletRequest request) throws UnsupportedEncodingException{

        Optional<SmtpVo> smtp = null;
        Claims claims = TokenUtils.getClaimsFormToken(request.getCookies());
        if(claims != null){
            smtp = smtpService.getSMTP(Long.parseLong(claims.get("user_no").toString()));
        }else{
            return new ResponseEntity<Object>("토큰을 확인해주세요", HttpStatus.FORBIDDEN);
        }

        try {
            smtp.get().setSmtpPw(aesService.decAES(smtp.get().getSmtpPw()));
            smtpService.connectCheck(smtp.get());
        }catch (IllegalStateException e) {
            e.printStackTrace();
            return new ResponseEntity<Object>("이미 연결이 되어있습니다.", HttpStatus.BAD_REQUEST);
        } catch(AuthenticationFailedException e){
            e.printStackTrace();
            return new ResponseEntity<Object>("로그인 정보가 올바르지 않습니다.", HttpStatus.UNAUTHORIZED);
        }catch (MessagingException e){
            e.printStackTrace();
            return new ResponseEntity<Object>("SMTP 설정을 다시 확인해주세요.", HttpStatus.BAD_REQUEST);
        }catch(Exception e){
            e.printStackTrace();
            return new ResponseEntity<Object>("SMTP 설정을 다시 확인해주세요.", HttpStatus.BAD_REQUEST);
        }

        return new ResponseEntity<Object>("성공", HttpStatus.OK);
    }


//    @ExceptionHandler(value = { MessagingException.class, NullPointerException.class, AuthenticationFailedException.class, NoSuchMethodError.class})
//    @ResponseStatus(value = HttpStatus.NOT_ACCEPTABLE)
//    public void nfeHandler(Exception e){
//        e.printStackTrace();
//    }

}

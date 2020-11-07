package com.hanium.mer.controller;

import com.hanium.mer.TokenUtils;
import com.hanium.mer.service.SMTPService;
import com.hanium.mer.vo.SmtpVo;
import com.hanium.mer.vo.Smtp_setting;
import io.jsonwebtoken.Claims;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import javax.servlet.http.HttpServletRequest;
import java.io.UnsupportedEncodingException;
import java.util.Optional;

//"		setting.GET(""/smtpSetting"", api.GetSmtpSetting)
//        setting.POST(""/smtpSetting"", api.SetSmtpSetting)
//        setting.POST(""/smtpConnectCheck"", api.SmtpConnectionCheck)"

@RestController
public class ApiSMTPController {

    @Autowired
    SMTPService smtpService;

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

    @PostMapping("/setting/smtpSetting")
    public HttpStatus setSTMPSetting(HttpServletRequest request, @RequestBody SmtpVo newSmtp) throws UnsupportedEncodingException {
        Optional<SmtpVo> smtp;
        Claims claims = TokenUtils.getClaimsFormToken(request.getCookies());
        if (claims != null) {
            try{
                smtpService.setSMTP(Long.parseLong(claims.get("user_no").toString()), newSmtp);
                System.out.println(newSmtp.toString());
                return HttpStatus.OK;
            }catch(Exception e){
                e.printStackTrace();
                //에러 400-> smtp 정보확인
                //401 비밀번호확인 제대로 설정하기
                return HttpStatus.BAD_REQUEST;
            }
        }

        return HttpStatus.FORBIDDEN;
    }

    @PostMapping("/setting/smtpConnectCheck")
    public HttpStatus connectSTMPTest(HttpServletRequest request, @RequestBody SmtpVo smtp) throws UnsupportedEncodingException {

        try {
            smtpService.connectCheck(smtp);
            return HttpStatus.OK;
        }catch (Exception e){
            e.printStackTrace();
            return HttpStatus.UNAUTHORIZED;
        }
    }
}

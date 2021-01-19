package com.hanium.mer.controller;

import com.hanium.mer.TokenUtils;
import com.hanium.mer.service.UserService;
import com.hanium.mer.vo.NewUserInfo;
import com.hanium.mer.vo.UserVo;
import com.hanium.mer.vo.User_info;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.ExpiredJwtException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.security.crypto.bcrypt.BCrypt;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.UnsupportedEncodingException;
import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

@Slf4j
@RestController
public class ApiUserSettingController {

    private static final int BLANK_ERROR = 0;
    private static final int NO_REGULAR = 1;
    private static final int SUCCESS = 2;

    @Autowired
    UserService userService;

    @GetMapping("/setting/userSetting")
    public ResponseEntity<Object> getUserSetting(HttpServletRequest request) throws UnsupportedEncodingException {
        Optional<UserVo> user;
        Claims claims = null;
        try {
            claims = TokenUtils.getClaimsFormToken(request.getCookies());
        }catch(ExpiredJwtException e){
            e.printStackTrace();
            return new ResponseEntity<Object>("error", HttpStatus.FORBIDDEN);
        }
        if(claims != null){
            user = userService.findByUserNo(Long.parseLong(claims.get("user_no").toString()));
            return new ResponseEntity<Object>(new User_info(user.get()), HttpStatus.OK);
        }

        return new ResponseEntity<Object>("error", HttpStatus.METHOD_NOT_ALLOWED);
    }


    //todo 현재비밀번호 확인시 해시코드로 확인?
    //todo userVO user_password 컬럼 지우기
    @PostMapping(value = "setting/userSetting") //consumes = MediaType.APPLICATION_FORM_URLENCODED_VALUE, user_no은 pw유효성 검사시만, service, repogitory 리팩터링해서 더 쉽게 암호얻어오고 수정
    public ResponseEntity<Object> setUserSetting(@RequestBody NewUserInfo newUserInfo, HttpServletRequest request, HttpServletResponse response) throws UnsupportedEncodingException{

        UserVo user;
        Claims claims = null;
        try {
            claims = TokenUtils.getClaimsFormToken(request.getCookies());
        }catch(ExpiredJwtException e){
            e.printStackTrace();
            return new ResponseEntity<Object>("token expired", HttpStatus.FORBIDDEN);
        }
        if(claims != null){
            user = userService.findByUserNo(Long.parseLong(claims.get("user_no").toString())).get();
            String currentPwHash = user.getUserPwHash();
            //서비스 로직나 dao로 옮기기

            if(!checkName(newUserInfo.getName()) || newUserInfo.getName().contains(" ") || newUserInfo.getName().length() == 0){
                return new ResponseEntity<Object>("올바른 이름을 입력하세요.", HttpStatus.BAD_REQUEST);
            }

            if(!BCrypt.checkpw(newUserInfo.getCurrent_pw(), currentPwHash)){
                return new ResponseEntity<Object>("비밀번호가 틀립니다.", HttpStatus.BAD_REQUEST); //400
            }

            //새로운 비밀번호 유효하지 않거나 check와 같지 않을 때
            if(!newUserInfo.getPassword().equals(newUserInfo.getPassword_check())){
                return new ResponseEntity<Object>( "new password or new password check error", HttpStatus.UNAUTHORIZED); //401
            }

            if(checkPassword(newUserInfo.getPassword()) == BLANK_ERROR){
                return new ResponseEntity<Object>("공백문자가 포함되어 있습니다.", HttpStatus.UNAUTHORIZED);
            }else if(checkPassword(newUserInfo.getPassword()) == NO_REGULAR){
                return new ResponseEntity<Object>("숫자, 영문자, 특수문자 하나씩 포함해주세요.", HttpStatus.UNAUTHORIZED);
            }

            userService.changeUserInfo(Long.parseLong(claims.get("user_no").toString()), newUserInfo);
            System.out.println(newUserInfo.toString());
            return new ResponseEntity<Object>("success", HttpStatus.OK); //200
        }

        return new ResponseEntity<Object>("error", HttpStatus.INTERNAL_SERVER_ERROR); //500
    }


    public boolean checkName(String name){
        String namePattern = "^[가-힣]*$";
        Matcher matcher = Pattern.compile(namePattern).matcher(name);
        if(matcher.matches()){
            return true;
        }
        return false;
    }

    public int checkPassword(String pw){
        String pwPattern = "^(?=.*\\d)(?=.*\\W)(?=.*[A-Za-z]).{8,16}$";
        Matcher matcher = Pattern.compile(pwPattern).matcher(pw);

        //같은 문자 4개 이상 사용불가
        //pwPattern = "(.)\\1\\1\\1";
        //Matcher sameCharacterMatcher = Pattern.compile(pwPattern).matcher(pw);

        //숫자,특수문자, 영 대소문자 조합 8~16자리일 경우
        if(!matcher.matches()){
            return NO_REGULAR;
        }

        if(pw.contains(" ")){
            return BLANK_ERROR;
        }

        return SUCCESS;
        /*

        if(sameCharacterMatcher.find()){

        }

        if(pw.contains(userId)){

        }

        */
    }

}

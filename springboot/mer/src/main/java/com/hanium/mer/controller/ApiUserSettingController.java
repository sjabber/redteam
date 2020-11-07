package com.hanium.mer.controller;

import com.hanium.mer.TokenUtils;
import com.hanium.mer.service.UserService;
import com.hanium.mer.vo.NewUserInfo;
import com.hanium.mer.vo.UserVo;
import com.hanium.mer.vo.User_info;
import io.jsonwebtoken.Claims;
import io.jsonwebtoken.ExpiredJwtException;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.Cookie;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.UnsupportedEncodingException;
import java.util.Optional;

@RestController
public class ApiUserSettingController {

    @Autowired
    UserService userService;

    //고쳐야 할것: responseEntity<적절한 객체>/ 에러/ 토큰서비스 직관적으로 수정/ header가 이렇게 오는게 맞나../ 쿠키가 아니라 헤더로 가져오던데../ cors, header 필터사용
    // 1. 토큰에서 얻은 email 부르기 => 토큰에서 claims값 가져오기
    // 2. 토큰에 있는 no으로 부르기 => 인덱스와 데이터 크기가 작기에 더 빠름
    // 3. 가져와서 수정이 아니라 바로 수정해줄 수도 있다.
    //@CrossOrigin(origins = "*", allowCredentials = "true", methods = {RequestMethod.GET, RequestMethod.POST, RequestMethod.PUT, RequestMethod.DELETE, RequestMethod.OPTIONS})
    @GetMapping("/setting/userSetting")
    public ResponseEntity<Object> getUserSetting(HttpServletRequest request, HttpServletResponse response) throws UnsupportedEncodingException {

        Optional<UserVo> user;
        Claims claims = null;
        try {
            claims = TokenUtils.getClaimsFormToken(request.getCookies());
        }catch(ExpiredJwtException e){
            e.printStackTrace();
            return new ResponseEntity<Object>("error", HttpStatus.FORBIDDEN);
        }
        if(claims != null){
            user = userService.findByUserNo(Long.parseLong(claims.get("user_no").toString()));;
            return new ResponseEntity<Object>(new User_info(user.get()), HttpStatus.OK);
        }

        return new ResponseEntity<Object>("error", HttpStatus.METHOD_NOT_ALLOWED);
    }


    @PostMapping(value = "setting/userSetting") //consumes = MediaType.APPLICATION_FORM_URLENCODED_VALUE, user_no은 pw유효성 검사시만, service, repogitory 리팩터링해서 더 쉽게 암호얻어오고 수정
    public HttpStatus setUserSetting(@RequestBody NewUserInfo newUserInfo,HttpServletRequest request, HttpServletResponse response) throws UnsupportedEncodingException{

        Optional<UserVo> user;
        Claims claims = null;
        try {
            claims = TokenUtils.getClaimsFormToken(request.getCookies());
        }catch(ExpiredJwtException e){
            e.printStackTrace();
            return HttpStatus.FORBIDDEN;
        }
        if(claims != null){
            userService.changeUserInfo(Long.parseLong(claims.get("user_no").toString()), newUserInfo);
            System.out.println(newUserInfo.toString());
            return HttpStatus.OK;
        }

        return HttpStatus.METHOD_NOT_ALLOWED;

        //안 될 경우, 특정에러와 함께 httpstatus/ 전체적인 restful exam 보고 로그랑 리팩터링 하기
        //postman은 raw에 json으로 해줘야하는데..흠 실제는 어떻게 될지..
    }


    //지워도됌, 확인용
    @PostMapping("/validatejwt")
    public String greet(HttpServletRequest request, HttpServletResponse response) throws UnsupportedEncodingException {
        request.getHeader("access-token");

        Cookie[] cookies = request.getCookies();
        Cookie cookie = null;

        for(Cookie c : cookies) {
            System.out.println("cookie name: " + c.getName());
            System.out.println("cookie value: " + c.getValue());

            if (c.getName().equals("access-token")) {
                cookie = c;
                if (TokenUtils.isValidToken(c.getValue()) == 0) {
                    return "hi";
                }
            }
        }

        return "unauthorized";
    }

    @GetMapping("/createjwt")
    public String getGreet(HttpServletRequest request, HttpServletResponse response) throws UnsupportedEncodingException {


        String token = TokenUtils.create();
        response.setHeader("access-token", token);
        System.out.println(token);

        return "create token";
    }
}

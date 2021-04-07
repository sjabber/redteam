package com.hanium.mer.interceptor;

import com.hanium.mer.TokenUtils;
import io.jsonwebtoken.ExpiredJwtException;
import io.jsonwebtoken.JwtException;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Component;
import org.springframework.web.servlet.HandlerInterceptor;
import org.springframework.web.servlet.ModelAndView;

import javax.servlet.http.Cookie;
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.text.SimpleDateFormat;
import java.util.Date;

@Slf4j
@Component
public class HttpHeaderNJWTInterceptor implements HandlerInterceptor{

    @Autowired
    TokenUtils tokenUtils;

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler) throws Exception
    {
        response = setHeader(response);

        Cookie[] cookies = request.getCookies();
        Cookie cookie = null;

        if (request.getMethod().equals("OPTIONS")){
            response.setStatus(HttpServletResponse.SC_OK);
            return true;
        }

        if (cookies != null) {
            for (Cookie c : cookies) {

                try{
                    if (c.getName().equals("access-token")) {
                        tokenUtils.isValidToken(c.getValue());
                        return true;
                    }
                } catch (ExpiredJwtException e) {

                    try {
                        response.sendError(403, "token expired");
                        log.info("token Expired");
                        return false;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                } catch (JwtException e) {

                    try {
                        response.sendError(405, "token tempered");
                        e.printStackTrace();
                        log.info("token tempered");
                        return true;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                } catch (NullPointerException e) {

                    try {
                        response.sendError(405, "token null");
                        log.info("toekn null");
                        return false;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                } catch(UnsupportedEncodingException e){

                    try {
                        response.sendError(405, "unsupportedEncoding");
                        log.info("unsupportedEncoding");
                        return false;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                }catch(NoSuchMethodError e){
                    try {
                        response.sendError(403, "토큰을 확인해주세요");
                        log.info("no such Method");
                        return false;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                }
            }
        }

        try {

            response.sendError(403, "null cookies");
            log.info("null cookies");
            return false;
        }catch(IOException ex){
            ex.printStackTrace();
            return false;
        }
    }

    @Override
    public void postHandle(HttpServletRequest request, HttpServletResponse response, Object handler, ModelAndView modelAndView)
            throws Exception {
    }

    @Override
    public void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex)
            throws Exception {
        // TODO Auto-generated method stub
    }

    public HttpServletResponse setHeader(HttpServletResponse response){
        SimpleDateFormat format = new SimpleDateFormat ( "yyyy-MM-dd HH:mm:ss");
        //cors설정에서 자동으로 붙여줌
        //하지만 에러시 controller까지 가지 않으므로 설정을 해줘야한다.
        //만약 preHandler에서 true라면 붙어서 cors헤더가 붙어서 올것이다.
        response.setHeader("Access-Control-Allow-Credentials", "true");
        //TODO 환경설정시 cors 포트 변경
        response.setHeader("Access-Control-Allow-Origin", "http://localhost:8080");
        response.setHeader("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE");
        response.setHeader("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept");
        response.setHeader("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0");
        response.setHeader("Last-Modified", format.format(new Date()));
        response.setHeader("Pragma", "no-cache");
        response.setHeader("Expires", "-1");

        return response;
    }
}

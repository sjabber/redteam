package com.hanium.mer.interceptor;

import com.hanium.mer.TokenUtils;
import io.jsonwebtoken.ExpiredJwtException;
import io.jsonwebtoken.JwtException;
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

@Component
public class HttpHeaderNJWTInterceptor implements HandlerInterceptor {

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler)
    {
        response = setHeader(response);

        Cookie[] cookies = request.getCookies();
        Cookie cookie = null;

        if (request.getMethod().equals("OPTIONS")){
            response.setStatus(HttpServletResponse.SC_OK);
            return true;
        }
        System.out.println("cookie num: "+ cookies.length);
        if (cookies != null) {
            for (Cookie c : cookies) {

                try{
                    if (c.getName().equals("access-token")) {
                        System.out.println(c.getValue());
                        TokenUtils.isValidToken(c.getValue());
                        return true;
                    }
                } catch (ExpiredJwtException e) {

                    try {
                        response.sendError(403, "token expired");
                        System.out.println("토큰 익스파이어");
                        return false;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                } catch (JwtException e) {

                    try {
                        response.sendError(405, "token tempered");
                        e.printStackTrace();
                        System.out.println("토큰 파괴");
                        return true;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                } catch (NullPointerException e) {

                    try {
                        response.sendError(405, "token null");
                        System.out.println("토큰 없음");
                        return false;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                } catch(UnsupportedEncodingException e){

                    try {
                        response.sendError(405, "unsupportedEncoding");
                        System.out.println("토큰 인코딩에러");
                        return false;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                }catch(NoSuchMethodError e){
                    try {
                        response.sendError(403, "토큰을 확인해주세요");
                        System.out.println("nosuchMetodh");
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
            System.out.println("쿠키없음");
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
        response.setHeader("Access-Control-Allow-Origin", "http://localhost:8888");
        response.setHeader("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE");
        response.setHeader("Access-Control-Allow-Headers","Origin, X-Requested-With, Content-Type, Accept");
        response.setHeader("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0");
        response.setHeader("Last-Modified", format.format(new Date()));
        response.setHeader("Pragma", "no-cache");
        response.setHeader("Expires", "-1");

        return response;
    }
}

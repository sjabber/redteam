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
    private static final int SUCCESS = 0;

    @Override
    public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler)
    {
        //최근까지 브라우저에서보면 헤더 붙지도 않음.. cors때문에 붙는거 제외하고는..cors에서 allow나, exponsed headers인가 해야하나..
        SimpleDateFormat format = new SimpleDateFormat ( "yyyy-MM-dd HH:mm:ss");
        //cors설정에서 자동으로 붙여줌
        response.setHeader("Access-Control-Allow-Credentials", "true");
        response.setHeader("Access-Control-Allow-Origin", "http://localhost:63342");
        response.setHeader("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE");
        response.setHeader("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0");
        response.setHeader("Last-Modified", format.format(new Date()));
        response.setHeader("Pragma", "no-cache");
        response.setHeader("Expires", "-1");

        Cookie[] cookies = request.getCookies();
        Cookie cookie = null;
        if (cookies != null) {
            for (Cookie c : cookies) {
                System.out.println("cookie exist");
                try{
                    if (c.getName().equals("access-token")) {
                        TokenUtils.isValidToken(c.getValue());
                        return true;
                    }
                } catch (ExpiredJwtException e) {
                    try {
                        response.sendError(405, "token expired");
                        return false;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                } catch (JwtException e) {
                    try {
                        response.sendError(405, "token tempered");
                        return false;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                } catch (NullPointerException e) {
                    try {
                        response.sendError(405, "token null");
                        return false;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                } catch(UnsupportedEncodingException e){
                    try {
                        response.sendError(405, "unsupportedEncoding");
                        return false;
                    }catch(IOException ex){
                        ex.printStackTrace();
                        return false;
                    }
                }
            }
        }

        try {
            response.sendError(444, "null cookies");
            return false;
        }catch(IOException ex){
            ex.printStackTrace();
            return false;
        }
    }

    @Override
    public void postHandle(HttpServletRequest request, HttpServletResponse response, Object handler, ModelAndView modelAndView)
            throws Exception {
//        SimpleDateFormat format = new SimpleDateFormat ( "yyyy-MM-dd HH:mm:ss");
//        //cors설정에서 자동으로 붙여줌
//        response.setHeader("Access-Control-Allow-Credentials", "true");
//        response.setHeader("Access-Control-Allow-Origin", "http://localhost:63342");
//        response.setHeader("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE");
//        response.setHeader("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0");
//        response.setHeader("Last-Modified", format.format(new Date()));
//        response.setHeader("Pragma", "no-cache");
//        response.setHeader("Expires", "-1");
    }

    @Override
    public void afterCompletion(HttpServletRequest request, HttpServletResponse response, Object handler, Exception ex)
            throws Exception {
        // TODO Auto-generated method stub
    }

    public HttpServletResponse setHeader(HttpServletResponse response){
        SimpleDateFormat format = new SimpleDateFormat ( "yyyy-MM-dd HH:mm:ss");
        //cors설정에서 자동으로 붙여줌
        //response.setHeader("Access-Control-Allow-Credentials", "true");
        //response.setHeader("Access-Control-Allow-Origin", "http://localhost:63342");
        //response.setHeader("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS, DELETE");
        response.setHeader("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0, max-age=0");
        response.setHeader("Last-Modified", format.format(new Date()));
        response.setHeader("Pragma", "no-cache");
        response.setHeader("Expires", "-1");

        return response;
    }
}

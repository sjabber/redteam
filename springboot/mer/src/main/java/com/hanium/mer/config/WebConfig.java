package com.hanium.mer.config;

import com.hanium.mer.interceptor.HttpHeaderNJWTInterceptor;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.context.annotation.Configuration;
import org.springframework.web.servlet.config.annotation.CorsRegistry;
import org.springframework.web.servlet.config.annotation.InterceptorRegistry;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

@Configuration
public class WebConfig implements WebMvcConfigurer {

    @Autowired
    HttpHeaderNJWTInterceptor httpHeaderNJWTInterceptor;



    //@CrossOrigin(origins = "*", allowCredentials = "true", methods = {RequestMethod.GET, RequestMethod.POST, RequestMethod.PUT, RequestMethod.DELETE, RequestMethod.OPTIONS})
    //Controller 메서드마다 붙이던 값을 전역으로 설정
    //postman에서는 안해도 되지만 실제 포트번호가 다르면 cors에러가 뜬다.
    // 헤더가 없어도 안된다.(에러는 제대로 보지 못했지만 userSetting 이동시 계속 로그인 페이지가 뜬다.), 될 때도 있다.. 브라우저캐시때문인가? 검은창으로해도그럼..
    @Override
    public void addCorsMappings(CorsRegistry registry) {
        registry.addMapping("/**")
                .allowedOrigins("*")
                .allowedMethods("*") //"GET", "POST","PUT", "DELETE"
                .allowCredentials(true);
    }


    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        registry.addInterceptor(httpHeaderNJWTInterceptor)
                .addPathPatterns("**");
    }

}

package com.hanium.mer.controller;

import com.hanium.mer.TokenUtils;
import com.hanium.mer.service.ProjectService;
import com.hanium.mer.vo.ProjectVo;
import io.jsonwebtoken.Claims;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import javax.servlet.http.HttpServletRequest;
import java.io.UnsupportedEncodingException;
import java.time.LocalDateTime;

@RestController
public class ApiProjectController {

    @Autowired
    ProjectService projectService;

    @PostMapping("/api/projectCreate")
    public ResponseEntity<Object> addProject(HttpServletRequest request, @RequestBody ProjectVo newProject)
            throws UnsupportedEncodingException {

        Claims claims = TokenUtils.getClaimsFormToken(request.getCookies());
        if (claims != null) {
            try{
                //todo 유효성검사
                newProject.setUserNo(Long.parseLong(claims.get("user_no").toString()));
                //JPA AUTO에 NULL값을 넣음. vo에서 따로 처리해줘도됌
                newProject.setCreatedTime(LocalDateTime.now());
                System.out.println(newProject.toString());
                projectService.addProject(newProject);
                return new ResponseEntity<Object>(newProject.toString(), HttpStatus.OK);
            }catch(Exception e){
                e.printStackTrace();
                return new ResponseEntity<Object>("프로젝트 생성 정보를 확인해주세요.", HttpStatus.BAD_REQUEST);
            }
        }

        return new ResponseEntity<Object>("토큰을 확인해주세요", HttpStatus.FORBIDDEN);
    }
}

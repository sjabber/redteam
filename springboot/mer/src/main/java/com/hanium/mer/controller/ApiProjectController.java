package com.hanium.mer.controller;

import com.hanium.mer.TokenUtils;
import com.hanium.mer.service.KafkaProducerService;
import com.hanium.mer.service.ProjectService;
import com.hanium.mer.service.SMTPService;
import com.hanium.mer.service.TemplateService;
import com.hanium.mer.vo.KafkaMessage;
import com.hanium.mer.vo.ProjectDto;
import com.hanium.mer.vo.ProjectVo;
import io.jsonwebtoken.Claims;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.util.StringUtils;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import javax.servlet.http.HttpServletRequest;
import java.io.UnsupportedEncodingException;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;

@Slf4j
@RestController
public class ApiProjectController {

    private final String NOTAG = "0";

    @Autowired
    ProjectService projectService;
    @Autowired
    SMTPService smtpService;
    @Autowired
    TemplateService templateService;

    @Autowired
    KafkaProducerService kafkaProducerService;



    @PostMapping("/api/projectCreate")
    public ResponseEntity<Object> addProject(HttpServletRequest request, @RequestBody ProjectDto projectDto)
            throws UnsupportedEncodingException {

        Claims claims = TokenUtils.getClaimsFormToken(request.getCookies());
        if (claims != null) {
            try{
                ProjectVo newProject = new ProjectVo();

                log.info("projectDto {}", projectDto);

                //todo 유효성검사(빈칸, 시작날짜, 끝 날짜 논리 순서), p_status는 시작 예약 종료? 한글로?
                newProject.setUserNo(Long.parseLong(claims.get("user_no").toString()));

                //JPA AUTO에 NULL값을 넣음. vo에서 따로 처리해줘도됌
                if(projectDto.getEnd_date().isBefore(projectDto.getStart_date())
                        || projectDto.getEnd_date().isBefore(LocalDate.now())){
                    return new ResponseEntity<Object>("날짜를 제대로 선택해주세요.", HttpStatus.BAD_REQUEST);
                }
                newProject.setStartDate(projectDto.getStart_date());
                newProject.setEndDate(projectDto.getEnd_date());

                if(StringUtils.isEmpty(projectDto.getP_name()) || StringUtils.isEmpty(projectDto.getP_description())){
                    return new ResponseEntity<Object>("프로젝트 정보를 제대로 생성해주세요.", HttpStatus.BAD_REQUEST);
                }
                newProject.setPName(projectDto.getP_name());
                newProject.setPDescription(projectDto.getP_description());

                if(projectDto.getTmp_no() == 0){
                    return new ResponseEntity<Object>("프로젝트 템플릿을 선택해주세요.", HttpStatus.BAD_REQUEST);
                }
                newProject.setTmlNo(projectDto.getTmp_no());

                //TODO 너무 비효율.. 중복제거와 first,second 대입 방법생각해보기
                if(projectDto.getTag_no().size() == 0){
                    return new ResponseEntity<Object>("프로젝트 태그가 하나 이상 선택해야합니다.", HttpStatus.BAD_REQUEST);
                }
                while(projectDto.getTag_no().size() < 3){
                    projectDto.getTag_no().add(NOTAG);
                }
                //태그 중복 불가, 0인 경우도 고려
                if((!projectDto.getTag_no().get(0).equals(NOTAG) && projectDto.getTag_no().get(0).equals(projectDto.getTag_no().get(1)))
                    ||  (!projectDto.getTag_no().get(2).equals(NOTAG)  && projectDto.getTag_no().get(0).equals(projectDto.getTag_no().get(2)))
                    || (!projectDto.getTag_no().get(1).equals(NOTAG) && projectDto.getTag_no().get(1).equals(projectDto.getTag_no().get(2)))){
                    return new ResponseEntity<Object>("프로젝트 태그가 중복됩니다.", HttpStatus.BAD_REQUEST);
                }
                newProject.setTagFirst(Integer.parseInt(projectDto.getTag_no().get(0)));
                newProject.setTagSecond(Integer.parseInt(projectDto.getTag_no().get(1)));
                newProject.setTagThird(Integer.parseInt(projectDto.getTag_no().get(2)));


                newProject.setCreatedTime(LocalDateTime.now());
                log.info("new project info: {}", newProject.toString());
                projectService.addProject(newProject);
                return new ResponseEntity<Object>(newProject.toString(), HttpStatus.OK);

            }catch(Exception e){
                e.printStackTrace();
                return new ResponseEntity<Object>("프로젝트 생성 정보를 확인해주세요.", HttpStatus.BAD_REQUEST);
            }
        }
        return new ResponseEntity<Object>("토큰을 확인해주세요", HttpStatus.FORBIDDEN);
    }

    @GetMapping("/api/getProjects")
    public ResponseEntity<Object> addProject(HttpServletRequest request)
            throws UnsupportedEncodingException {

        Claims claims = TokenUtils.getClaimsFormToken(request.getCookies());
        if (claims != null) {
            try {
                Map<String,Object> map = new HashMap<>();
                List<ProjectVo> project_list = projectService.getProjects(Long.parseLong(claims.get("user_no").toString()));
                map.put("projects", project_list);
                return new ResponseEntity<Object>(map, HttpStatus.OK);
            }catch(Exception e){
                e.printStackTrace();
                return new ResponseEntity<Object>("DB 정보를 확인해주세요.", HttpStatus.BAD_REQUEST);
            }
        }
        return new ResponseEntity<Object>("토큰을 확인해주세요", HttpStatus.FORBIDDEN);
    }

    @GetMapping("/api/kafka")
    public ResponseEntity<Object> sendMessageTest(HttpServletRequest request, Long p_no)
            throws UnsupportedEncodingException {

        //TODO token user_id 가져오기
        Claims claims = TokenUtils.getClaimsFormToken(request.getCookies());
        System.out.println(TokenUtils.create());
        if (claims != null) {
            try {
                //TODO 프로젝트 번호로 받아오기
                Optional<ProjectVo> project = projectService.getProject(p_no, Long.parseLong(claims.get("user_no").toString()));
                //log.info(project.get().toString());

                List<Object[]> targets = projectService.getTargets(project.get().getUserNo(), project.get().getTagFirst(),
                        project.get().getTagSecond(), project.get().getTagThird()); //line 70
                //log.info("project target query result {}", targets.toString());

                for (Object[] m : targets) {
                    KafkaMessage kafkaMessage = new KafkaMessage();
                    kafkaMessage.setTargetNo( (int) m[0]);
                    kafkaMessage.setPNo(project.get().getPNo());
                    kafkaMessage.setTmpNo(project.get().getTmlNo());
                    kafkaMessage.setUserNo(project.get().getUserNo());

                    kafkaProducerService.sendMessage(kafkaMessage);
                }

                return new ResponseEntity<Object>("success", HttpStatus.OK);
            }catch(Exception e){
                e.printStackTrace();
                return new ResponseEntity<Object>("카프카 정보를 확인해주세요.", HttpStatus.BAD_REQUEST);
            }
        }
        return new ResponseEntity<Object>("토큰을 확인해주세요", HttpStatus.OK);

    }

}

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
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import javax.servlet.http.HttpServletRequest;
import java.io.UnsupportedEncodingException;
import java.time.LocalDate;
import java.time.LocalDateTime;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

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

                String pattern = "^[a-zA-Z0-9가-힣]{2,100}$";
                Matcher matcher = Pattern.compile(pattern).matcher(projectDto.getP_name());
                if(!matcher.matches()){
                    return new ResponseEntity<Object>("프로젝트 이름을 확인해주세요.", HttpStatus.BAD_REQUEST);
                }
                newProject.setPName(projectDto.getP_name());

                if(projectDto.getP_description().length() >= 1000){
                    return new ResponseEntity<Object>("프로젝트 설명을 확인해주세요.", HttpStatus.BAD_REQUEST);
                }
                newProject.setPDescription(projectDto.getP_description());

                //JPA AUTO에 NULL값을 넣음. vo에서 따로 처리해줘도됌
                //log.info("{} {} {}", projectDto.getStart_date(), LocalDate.now(), projectDto.getStart_date().isBefore(LocalDate.now()));
                if( projectDto.getEnd_date().isBefore(projectDto.getStart_date())
                        || projectDto.getEnd_date().isBefore(LocalDate.now())
                        || projectDto.getStart_date().isBefore(LocalDate.now())){
                    log.info("{} {} {}", projectDto.getStart_date(), LocalDate.now(), projectDto.getStart_date().isBefore(LocalDate.now()));
                    return new ResponseEntity<Object>("날짜를 제대로 선택해주세요.", HttpStatus.BAD_REQUEST);
                }
                newProject.setStartDate(projectDto.getStart_date());
                newProject.setEndDate(projectDto.getEnd_date());

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

                //TODO p_status 설정, 안해도 자동 진행으로 저장됨 db default값
                if(LocalDate.now().isBefore(projectDto.getStart_date())){
                    newProject.setPStatus(2);
                }else{
                    newProject.setPStatus(1);
                }
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

    //TODO page, size를 front에서 파라미터로 넘기면 됨, 정렬도 생각해보기
    //TODO front 코드 가져오면 됨
    @GetMapping("/api/getProjects")
    public ResponseEntity<Object> addProject(HttpServletRequest request, @RequestParam(defaultValue = "1") int page)
            throws UnsupportedEncodingException {

        Claims claims = TokenUtils.getClaimsFormToken(request.getCookies());
        if (claims != null) {
            try {
                Map<String,Object> map = new HashMap<>();
                Pageable pageable = PageRequest.of(page - 1, 5, Sort.by("createdTime").descending());
                Page<ProjectVo> project_list = projectService.getProjects(Long.parseLong(claims.get("user_no").toString()), pageable);
                map.put("project_list", project_list);
                log.info(map.toString());
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
                    kafkaMessage.setTarget_no( (int) m[0]);
                    kafkaMessage.setP_no(project.get().getPNo());
                    kafkaMessage.setTmp_no(project.get().getTmlNo());
                    kafkaMessage.setUser_no(project.get().getUserNo());

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

package com.hanium.mer.service;


import com.fasterxml.jackson.databind.ObjectMapper;
import com.hanium.mer.repogitory.TargetRepository;
import com.hanium.mer.vo.KafkaMessage;
import com.hanium.mer.vo.TargetVo;
import com.hanium.mer.vo.TemplateVO;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;

import java.io.IOException;
import java.util.Optional;

@Slf4j
@Service
public class KafkaConsumerService {
    private final String COUNT_ADDRESS = "http://localhost:5000/api/CountTarget?";
    private final String FISHING_ADDRESS = "http://localhost:8888/warn/google2?";

    @Autowired
    SMTPService smtpService;
    @Autowired
    TargetRepository targetRepository;
    @Autowired
    ProjectService projectService;

    //TODO groupId 수정 필요
    @KafkaListener(topics = "redteam", groupId = "RedTeam")
    public void consume(String strMessage) throws IOException {
        ObjectMapper objectMapper =new ObjectMapper();
        KafkaMessage message = objectMapper.readValue(strMessage, KafkaMessage.class);
        log.info(String.format("Consumed message : %s", message));
        try {
            //TODO jpa 조인 사용해보기
            Optional<TargetVo> target = targetRepository.findByTargetNo(message.getTarget_no());
            Optional<TemplateVO> template = projectService.getTemplate(message.getTmp_no());
            //TODO replaceAll써야할수도, link ip 추가 시 수정 필요
            String mailContent = template.get().getMailContent();
            if(target.isPresent()) {
                mailContent = mailContent.replace("{{target_name}}", target.get().getTargetName());
                mailContent = mailContent.replace("{{target_position}}", target.get().getTargetPosition());
                mailContent = mailContent.replace("{{target_organize}}", target.get().getTargetOrganize());
                mailContent = mailContent.replace("{{target_phone}}", target.get().getTargetPhone());
                //String linkIpStr = "<a href='#' onclick =\'function callFunction(){const r = new XMLHttpRequest(); r.open(\"GET\", \"http://localhost:5000/api/CountTarget?tNo="+ target.get().getTargetNo() + "&pNo="+ message.getP_no() +"&email=false&link=true&download=false\", true);r.setRequestHeader(\"Content-Type\", \"application/x-www-form-urlencoded\"); r.withCredentials = true; r.onreadystatechange = function () { if (r.readyState === 4) { if (r.status === 200) {  document.location.href = \"http://localhost:8888/warn/warning2\"; } else {document.location.reload()}    } }; r.send(); }; return false;'> 내용</a>";
                mailContent = mailContent.replace("{{link_ip}}", "<a href='http://localhost:5001/warning?tNo="+message.getTarget_no()+"&pNo="+message.getP_no()+"'>홈으로</a>");
                mailContent = mailContent.replace("{{count_ip}}", "<img src='"+COUNT_ADDRESS + "tNo="+target.get().getTargetNo()+
                        "&pNo="+message.getP_no()+"&email=true&link=false&download=false'"+" style='height: 0px; width: 0px'>");
                mailContent = mailContent.replace("{{fishing_ip}}", "<a href='"+FISHING_ADDRESS + "tNo="+target.get().getTargetNo()+
                        "&pNo="+message.getP_no()+"&email=false&link=true&download=false'>계정 보호하기 </a>");
            }

            if(template.get().getDownloadType() == 2){
                /*
                mailContent += "<a href='http://localhost:5001/api/file?fileName=" +
                        message.getP_no() +  "_" + message.getTarget_no() + "_bat"
                        + "'>파일 다운로드</a>";
                */
                mailContent += "<a href='http://localhost:5001/api/file?tNo=" +
                        message.getTarget_no() + "&pNo=" + message.getP_no() + "'>파일 다운로드</a>";
            }
            //System.out.println(mailContent);
            template.get().setMailContent(mailContent);

            log.info("consumer target {}", target.get());
            //log.info("consumer template {}", template.get());

            smtpService.sendMail(smtpService.getSMTP(message.getUser_no()).get(), target.get(), template.get(), message.getP_no());
            //smtpService.sendMail(message.getSmtp(), message.getTarget(), message.getTemplate());

            //TODO sendTo업데이트 하기
            projectService.updateSendNo(projectService.getProject(message.getP_no(), message.getUser_no()).get());
        }catch(Exception e){
            e.printStackTrace();
        }
    }

//    //TODO groupId 수정 필요
//    @KafkaListener(topics = "redteam", groupId = "RedTeam")
//    public void consume(KafkaMessage message) throws IOException {
//        log.info(String.format("Consumed message : %s", message));
//        try {
//        //TODO jpa 조인 사용해보기
//            Optional<TargetVo> target = targetRepository.findByTargetNo(message.getTarget_no());
//            Optional<TemplateVO> template = projectService.getTemplate(message.getTmp_no());
//            //TODO replaceAll써야할수도, link ip 추가 시 수정 필요
//            String mailContent = template.get().getMailContent();
//            if(target.isPresent()) {
//                mailContent = mailContent.replace("{{target_name}}", target.get().getTargetName());
//                mailContent = mailContent.replace("{{target_position}}", target.get().getTargetPosition());
//                mailContent = mailContent.replace("{{target_organize}}", target.get().getTargetOrganize());
//                mailContent = mailContent.replace("{{target_phone}}", target.get().getTargetPhone());
//                //String linkIpStr = "<a href='#' onclick =\'function callFunction(){const r = new XMLHttpRequest(); r.open(\"GET\", \"http://localhost:5000/api/CountTarget?tNo="+ target.get().getTargetNo() + "&pNo="+ message.getP_no() +"&email=false&link=true&download=false\", true);r.setRequestHeader(\"Content-Type\", \"application/x-www-form-urlencoded\"); r.withCredentials = true; r.onreadystatechange = function () { if (r.readyState === 4) { if (r.status === 200) {  document.location.href = \"http://localhost:8888/warn/warning2\"; } else {document.location.reload()}    } }; r.send(); }; return false;'> 내용</a>";
//                mailContent = mailContent.replace("{{link_ip}}", "<a href='http://localhost:5001/warning?tNo="+message.getTarget_no()+"&pNo="+message.getP_no()+"'>홈으로</a>");
//                mailContent = mailContent.replace("{{count_ip}}", "<img src='"+COUNT_ADDRESS + "tNo="+target.get().getTargetNo()+
//                        "&pNo="+message.getP_no()+"&email=true&link=false&download=false'"+" style='height: 0px; width: 0px'>");
//            }
//            System.out.println(mailContent);
//            template.get().setMailContent(mailContent);
//
//            log.info("consumer target {}", target.get());
//            //log.info("consumer template {}", template.get());
//
//            smtpService.sendMail(smtpService.getSMTP(message.getUser_no()).get(), target.get(), template.get(), message.getP_no());
//            //smtpService.sendMail(message.getSmtp(), message.getTarget(), message.getTemplate());
//
//            //TODO sendTo업데이트 하기
//            projectService.updateSendNo(projectService.getProject(message.getP_no(), message.getUser_no()).get());
//        }catch(Exception e){
//            e.printStackTrace();
//        }
//    }

}

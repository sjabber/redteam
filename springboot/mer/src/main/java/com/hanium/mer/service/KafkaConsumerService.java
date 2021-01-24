package com.hanium.mer.service;


import lombok.extern.slf4j.Slf4j;
import org.springframework.stereotype.Service;

@Slf4j
@Service
public class KafkaConsumerService {
/*
    private final String COUNT_ADDRESS = "http://localhost:5000/api/CountTarget?";

    @Autowired
    SMTPService smtpService;
    @Autowired
    TargetRepository targetRepository;
    @Autowired
    ProjectService projectService;

    //TODO groupId 수정 필요
    @KafkaListener(topics = "redteam", groupId = "RED")
    public void consume(KafkaMessage message) throws IOException {
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
                mailContent = mailContent.replace("{{count_ip}}", "<img src="+COUNT_ADDRESS + "tNo="+target.get().getTargetNo()+
                        "&pNo="+message.getP_no()+"&email=true&link=false&download=True >");
            }
            template.get().setMailContent(mailContent);

            log.info("consumer target {}", target.get());
            //log.info("consumer template {}", template.get());

            smtpService.sendMail(smtpService.getSMTP(message.getUser_no()).get(), target.get(), template.get());
            //smtpService.sendMail(message.getSmtp(), message.getTarget(), message.getTemplate());

            //TODO sendTo 업데이트 하기
        }catch(Exception e){
            e.printStackTrace();
        }
    }
*/
}

package com.hanium.mer.service;

import com.hanium.mer.repogitory.ProjectRepository;
import com.hanium.mer.repogitory.TemplateRepository;
import com.hanium.mer.vo.ProjectVo;
import com.hanium.mer.vo.TemplateVO;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Optional;

@Slf4j
@Service
public class ProjectService {

    @Autowired
    ProjectRepository projectRepository;

    @Autowired
    TemplateRepository templateRepository;

    public Optional<ProjectVo> getProject(Long p_no, Long user_no){
        log.info("get Project p_no {}", p_no);
        return projectRepository.findBypNoAndUserNo(p_no, user_no);
    }

    public Optional<TemplateVO> getTemplate(Long tmpNo){
        log.info("get tmpNo {}", tmpNo);
        return templateRepository.findByTmpNo(tmpNo);
    }

    public void updateSendNo(ProjectVo projectVo){
        projectVo.setSendNo(projectVo.getSendNo()+1);
        projectRepository.save(projectVo);
    }

}

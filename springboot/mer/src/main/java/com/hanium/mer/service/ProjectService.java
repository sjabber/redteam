package com.hanium.mer.service;

import com.hanium.mer.repogitory.ProjectRepository;
import com.hanium.mer.repogitory.TemplateRepository;
import com.hanium.mer.vo.ProjectVo;
import com.hanium.mer.vo.TemplateVO;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Slf4j
@Service
public class ProjectService {

    @Autowired
    ProjectRepository projectRepository;

    @Autowired
    TemplateRepository templateRepository;

    public void addProject(ProjectVo newProject){
        log.info("newProject");
        projectRepository.save(newProject);
    }

    public Optional<ProjectVo> getProject(Long p_no, Long user_no){
        log.info("get Project p_no {}", p_no);
        return projectRepository.findBypNoAndUserNo(p_no, user_no);
    }

    public Page<ProjectVo> getProjects(Long user_no, Pageable pageable){
        //log.info("get Projects user_no {}", user_no);
        return projectRepository.findByUserNo(user_no, pageable);
    }

    public List<Object[]> getTargets(Long userNo, int tagFirst, int tagSecond, int tagThird){
        log.info(projectRepository.findTargetByUserNoAndTags(userNo, tagFirst, tagSecond, tagThird).toString());
        return projectRepository.findTargetByUserNoAndTags(userNo, tagFirst, tagSecond, tagThird);
    }

    public Optional<TemplateVO> getTemplate(Long tmpNo){
        log.info("get tmpNo {}", tmpNo);
        return templateRepository.findByTmpNo(tmpNo);
    }
}

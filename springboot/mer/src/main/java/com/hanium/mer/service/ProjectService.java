package com.hanium.mer.service;

import com.hanium.mer.repogitory.ProjectRepository;
import com.hanium.mer.vo.ProjectVo;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

@Service
public class ProjectService {

    @Autowired
    ProjectRepository projectRepository;

    public void addProject(ProjectVo newProject){
        System.out.println(newProject.toString());
        projectRepository.save(newProject);
    }
}

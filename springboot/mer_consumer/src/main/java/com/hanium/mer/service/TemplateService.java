package com.hanium.mer.service;

import com.hanium.mer.repogitory.TemplateRepository;
import com.hanium.mer.vo.TemplateVO;
import lombok.extern.slf4j.Slf4j;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.Optional;

@Slf4j
@Service
public class TemplateService {

    @Autowired
    TemplateRepository templateRepository;

    public Optional<TemplateVO> getTemplate(Long tmp_no){
        return templateRepository.findByTmpNo(tmp_no);
    }

}

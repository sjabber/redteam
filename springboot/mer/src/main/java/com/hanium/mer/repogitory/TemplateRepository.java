package com.hanium.mer.repogitory;

import com.hanium.mer.vo.TemplateVO;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;


@Repository
public interface TemplateRepository extends JpaRepository<TemplateVO, Long> {

    //TODO userNo 추가
    Optional<TemplateVO> findByTmpNo(Long tmp_no);

}

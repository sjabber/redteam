package com.hanium.mer.repogitory;

import com.hanium.mer.vo.ProjectVo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;


@Repository
public interface ProjectRepository extends JpaRepository<ProjectVo, Long> {
    Optional<ProjectVo> findBypNoAndUserNo(Long p_no, Long user_no);

}

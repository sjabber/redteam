package com.hanium.mer.repogitory;

import com.hanium.mer.vo.ProjectVo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;


@Repository
public interface ProjectRepository extends JpaRepository<ProjectVo, Long> {

}

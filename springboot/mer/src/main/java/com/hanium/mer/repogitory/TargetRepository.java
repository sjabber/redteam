package com.hanium.mer.repogitory;


import com.hanium.mer.vo.TargetVo;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface TargetRepository extends JpaRepository<TargetVo, Long> {

    public Optional<TargetVo> findByTargetNo(int target_no);
}

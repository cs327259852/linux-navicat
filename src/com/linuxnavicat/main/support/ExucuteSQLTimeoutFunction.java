package com.linuxnavicat.main.support;

import com.linuxnavicat.main.exception.SQLBadGrammarException;

@FunctionalInterface
public interface ExucuteSQLTimeoutFunction<T,R> {

    /**
     *
     * @param t
     * @return
     * @throws SQLBadGrammarException
     */
    R apply(T t)throws SQLBadGrammarException;
}

package com.linuxnavicat.main.listener;

import com.linuxnavicat.main.entitys.MessageLocale;

/**
 * 系统语言监听者接口
 */
@FunctionalInterface
public interface SystemLanguageListener  {

    void systemLanguageChanged(MessageLocale messageLocale);
}

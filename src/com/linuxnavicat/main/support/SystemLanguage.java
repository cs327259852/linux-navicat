package com.linuxnavicat.main.support;

import com.linuxnavicat.main.entitys.MessageLocale;
import com.linuxnavicat.main.listener.SystemLanguageListener;

import java.util.ArrayList;
import java.util.List;

/**
 * 系统语言事件对象
 */
public class SystemLanguage {

    /**
     * 系统语言事件监听者
     */
    private List<SystemLanguageListener> listenerList = new ArrayList<>();

    public boolean addListener(SystemLanguageListener listener){
        return listenerList.add(listener);
    }

    public void systemLanguageNotify(MessageLocale messageLocale){
        for(SystemLanguageListener listener:listenerList){
            listener.systemLanguageChanged(messageLocale);
        }
    }
}

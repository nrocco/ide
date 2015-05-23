<?php

// This is a hack to avoid having local modifications in all app_dev.php front controllers
if ("/app_dev.php" === $_SERVER['SCRIPT_NAME']) {
    $_SERVER['REMOTE_ADDR'] = '127.0.0.1';
}

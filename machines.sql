-- phpMyAdmin SQL Dump
-- version 3.4.8
-- http://www.phpmyadmin.net
--
-- Хост: localhost
-- Время создания: Сен 08 2015 г., 15:52
-- Версия сервера: 5.0.67
-- Версия PHP: 5.3.6

SET SQL_MODE="NO_AUTO_VALUE_ON_ZERO";
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;

--
-- База данных: `machines_data`
--

-- --------------------------------------------------------

--
-- Структура таблицы `machines`
--

CREATE TABLE IF NOT EXISTS `machines` (
  `id` int(11) NOT NULL auto_increment,
  `ip` varchar(15) NOT NULL,
  `table_name` varchar(25) NOT NULL,
  `unique_id` int(11) default NULL COMMENT 'machine''s id',
  `title` varchar(30) default NULL,
  PRIMARY KEY  (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=14 ;

--
-- Дамп данных таблицы `machines`
--

INSERT INTO `machines` (`id`, `ip`, `table_name`, `unique_id`, `title`) VALUES
(1, '_140.1.24.79', 'raute_shell', 1, NULL),
(2, '_140.1.24.81', 'raute_cutter', 2, NULL),
(3, '140.1.27.100', 'boiler', NULL, NULL),
(4, '_140.1.24.220', 'warmer', 4, NULL),
(5, '140.1.25.100', 'lopping', NULL, NULL),
(6, '140.1.26.100', 'grinding', NULL, NULL),
(7, '10.1.190.222', 'merger', NULL, NULL),
(8, '_140.1.24.15', 'saw', 8, NULL),
(9, '10.1.191.222', 'merger_old', NULL, NULL),
(10, '140.140.140.140', 'press', 10, NULL),
(11, '_192.168.0.231', 'warmer_out', 11, NULL),
(12, '_0.0.0.0', 'saw_new', 12, 'Новая распиловка'),
(13, '_140.1.24.51', 'edge_gluing', 13, 'Ребросклейка');

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

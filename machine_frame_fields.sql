-- phpMyAdmin SQL Dump
-- version 3.4.8
-- http://www.phpmyadmin.net
--
-- Хост: localhost
-- Время создания: Сен 08 2015 г., 15:53
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
-- Структура таблицы `machine_frame_fields`
--

CREATE TABLE IF NOT EXISTS `machine_frame_fields` (
  `id` int(11) NOT NULL auto_increment,
  `machine_id` int(11) NOT NULL COMMENT 'ID станка',
  `field_index` int(11) NOT NULL COMMENT 'Порядковый номер поля в кадре',
  `field_name` varchar(25) character set latin1 NOT NULL COMMENT 'Имя поля для таблицы',
  `field_size` int(11) NOT NULL COMMENT 'Ширина боля в байтах',
  `field_type` varchar(25) character set latin1 NOT NULL COMMENT 'Название типа для SQL запроса',
  `field_title` varchar(45) NOT NULL,
  PRIMARY KEY  (`id`)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 AUTO_INCREMENT=122 ;

--
-- Дамп данных таблицы `machine_frame_fields`
--

INSERT INTO `machine_frame_fields` (`id`, `machine_id`, `field_index`, `field_name`, `field_size`, `field_type`, `field_title`) VALUES
(2, 1, 0, 'r_init', 2, 'INT', 'Радиус поступившего чурака, мм'),
(3, 1, 1, 'r_cyl', 2, 'INT', 'Радиус цилиндрованного чурака, мм'),
(4, 1, 2, 'conisity', 2, 'INT', 'Конусность, мм/м'),
(5, 1, 3, 'r_core', 2, 'INT', 'Радиус карандаша, мм'),
(6, 1, 4, 'length', 2, 'INT', 'Длина, мм'),
(11, 1, 5, 'percent_cyl', 2, 'INT', '% на цилиндровку'),
(12, 1, 6, 'percent_core', 2, 'INT', '% на карандаш'),
(13, 1, 7, 'percent_exit', 2, 'INT', '% выхода'),
(14, 1, 8, 'curr_vol', 4, 'INT', 'Объём чурака, мм<sup>3</sup>'),
(15, 2, 0, 'thickness', 2, 'INT', 'Толщина, мм'),
(16, 2, 1, 'length', 2, 'INT', 'Длина, мм'),
(17, 2, 2, 'bin1_width', 2, 'INT', 'Ширина на столе 1, мм'),
(18, 2, 3, 'bin2_width', 2, 'INT', 'Ширина на столе 2, мм'),
(19, 2, 4, 'bin3_width', 2, 'INT', 'Ширина на столе 3, мм'),
(20, 2, 5, 'bin1_incr', 4, 'INT', 'Прирост объёма на столе 1, мм<sup>3</sup>'),
(21, 2, 6, 'bin2_incr', 4, 'INT', 'Прирост объёма на столе 2, мм<sup>3</sup>'),
(22, 2, 7, 'bin3_incr', 4, 'INT', 'Прирост объёма на столе 3, мм<sup>3</sup>'),
(26, 2, 8, 'bin1_quant', 2, 'INT', 'Прирост листов на столе 1'),
(27, 2, 9, 'bin2_quant', 2, 'INT', 'Прирост листов на столе 2'),
(28, 2, 10, 'bin3_quant', 2, 'INT', 'Прирост листов на столе 3'),
(29, 1, 9, 'core_vol', 4, 'INT', 'Объём сброшенного карандаша, мм<sup>3</sup>'),
(30, 3, 0, 'temp_specified', 2, 'INT', 'Заданная температура'),
(31, 3, 1, 'temp_send', 2, 'INT', 'Прямая температура'),
(32, 3, 2, 'temp_receive', 2, 'INT', 'Обратная температура'),
(33, 3, 3, 'power', 2, 'INT', 'Мощность'),
(34, 4, 0, 'packages', 2, 'INT', 'Число партий в минуту'),
(35, 4, 1, 'velocity', 2, 'INT', 'Скорость внутри сушилки'),
(40, 4, 2, 'totalSq', 4, 'INT', 'Поданная площадь, м²'),
(41, 5, 0, 'part_id', 2, 'INT', 'ID партии'),
(42, 5, 1, 'length', 2, 'INT', 'Длина'),
(43, 5, 2, 'width', 2, 'INT', 'Ширина'),
(44, 5, 3, 'velocity_short', 2, 'INT', 'Скорость на короткой опиловке'),
(45, 5, 4, 'velocity_long', 2, 'INT', 'Скорость на длинной опиловке'),
(46, 5, 5, 'value', 4, 'INT', 'Объём, мм^3'),
(47, 5, 6, 'thickness', 2, 'INT', 'Толщина, мм'),
(48, 6, 0, 'part_id', 2, 'INT', 'ID листа'),
(49, 6, 1, 'length', 2, 'INT', 'Длина листа'),
(50, 6, 2, 'width', 2, 'INT', 'Ширина, мм'),
(51, 6, 3, 'thickness', 2, 'INT', 'Толщина, мм'),
(52, 6, 4, 'value_bin1', 4, 'INT', 'Объем на первый стол, мм3'),
(53, 6, 5, 'value_bin2', 4, 'INT', 'Объем на второй стол, мм3'),
(54, 6, 6, 'value_bin3', 4, 'INT', 'Объем на третий стол, мм3'),
(55, 6, 7, 'velocity', 2, 'INT', 'Скорость шлифовки, м/мин'),
(56, 7, 0, 'id_send', 2, 'INT', 'id отправляемой партии'),
(57, 7, 1, 'count', 2, 'INT', 'количество пройденных листов'),
(58, 7, 2, 'length', 2, 'INT', 'длина опиловки, мм'),
(59, 7, 3, 'width', 2, 'INT', 'ширина опиловки, мм'),
(60, 7, 4, 'press1_merg_count', 2, 'INT', 'количество сращиваний 1'),
(61, 7, 5, 'press1_cut_count', 2, 'INT', 'количество отрубаний 1'),
(62, 7, 6, 'press1_length', 2, 'INT', 'длина готового листа 1'),
(63, 7, 7, 'press2_merg_count', 2, 'INT', 'количество сращиваний 2'),
(64, 7, 8, 'press2_cut_count', 2, 'INT', 'количество отрубаний 2'),
(65, 7, 9, 'press2_length', 2, 'INT', 'длина готового листа 2'),
(66, 7, 10, 'press3_merg_count', 2, 'INT', 'количество сращиваний 3'),
(67, 7, 11, 'press3_cut_count', 2, 'INT', 'количество отрубаний 3'),
(68, 7, 12, 'press3_length', 2, 'INT', 'длина готового листа 3'),
(69, 7, 13, 'press1_merg_value', 4, 'INT', 'Прирост объёма на прессе 1'),
(70, 7, 14, 'press2_merg_value', 4, 'INT', 'Прирост объёма на прессе 2'),
(71, 7, 15, 'press3_merg_value', 4, 'INT', 'Прирост объёма на прессе 3'),
(72, 7, 16, 'merg_total_value', 4, 'INT', 'Суммарный прирост объёма'),
(73, 8, 0, 'tree_format', 2, 'INT', 'Формат поступающих стволов'),
(74, 9, 0, 'id_send', 2, 'INT', 'id отправляемой партии'),
(75, 9, 1, 'count', 2, 'INT', 'количество пройденных листов'),
(76, 9, 2, 'length', 2, 'INT', 'длина опиловки, мм'),
(77, 9, 3, 'width', 2, 'INT', 'ширина опиловки, мм'),
(78, 9, 4, 'press1_merg_count', 2, 'INT', 'количество сращиваний 1'),
(79, 9, 5, 'press1_cut_count', 2, 'INT', 'количество отрубаний 1'),
(80, 9, 6, 'press1_length', 2, 'INT', 'длина готового листа 1'),
(81, 9, 7, 'press2_merg_count', 2, 'INT', 'количество сращиваний 2'),
(82, 9, 8, 'press2_cut_count', 2, 'INT', 'количество отрубаний 2'),
(83, 9, 9, 'press2_length', 2, 'INT', 'длина готового листа 2'),
(84, 9, 10, 'press3_merg_count', 2, 'INT', 'количество сращиваний 3'),
(85, 9, 11, 'press3_cut_count', 2, 'INT', 'количество отрубаний 3'),
(86, 9, 12, 'press3_length', 2, 'INT', 'длина готового листа 3'),
(87, 9, 13, 'press1_merg_value', 4, 'INT', 'Прирост объёма на прессе 1'),
(88, 9, 14, 'press2_merg_value', 4, 'INT', 'Прирост объёма на прессе 2'),
(89, 9, 15, 'press3_merg_value', 4, 'INT', 'Прирост объёма на прессе 3'),
(90, 9, 16, 'merg_total_value', 4, 'INT', 'Суммарный прирост объёма'),
(91, 10, 0, 'id_send', 2, 'INT', 'ID пакета (циклически)'),
(92, 10, 1, 'pressure', 2, 'INT', 'Давление'),
(93, 10, 2, 'temp', 2, 'INT', 'Температура'),
(94, 10, 3, 'width', 2, 'INT', 'Ширина'),
(95, 10, 4, 'length', 2, 'INT', 'Длина'),
(96, 10, 5, 'thickness', 2, 'INT', 'Толщина'),
(97, 10, 6, 'quant', 2, 'INT', 'Количество'),
(98, 10, 7, 'value', 4, 'INT', 'Объём'),
(99, 10, 8, 'time_1', 2, 'INT', 'Время 1'),
(100, 10, 9, 'time_2', 2, 'INT', 'Время 2'),
(101, 11, 0, 'usefull_square', 4, 'INT', 'Полезная площадь, сотые м²'),
(102, 11, 1, 'usefull_quant', 2, 'INT', 'Количество полезных листов, шт'),
(103, 11, 4, 'useless_square', 4, 'INT', 'Площадь отходов, сотые м²'),
(104, 11, 5, 'useless_quant', 2, 'INT', 'Количество листов на отходы, шт'),
(105, 11, 2, 'merger_square', 4, 'INT', 'Площадь листов на ребросклейку'),
(106, 11, 3, 'merger_quant', 2, 'INT', 'Количество листов на ребросклейку'),
(107, 4, 3, 'redSquare', 4, 'INT', 'Площадь красного шпона'),
(108, 4, 4, 'shelves', 2, 'INT', 'Количество полок'),
(109, 12, 0, 'id_send', 2, 'INT', 'Counter'),
(110, 12, 1, 'init_length', 2, 'INT', 'Длина ствола'),
(111, 12, 2, 'peeler_length', 2, 'INT', 'Установленная длина чурака'),
(112, 12, 3, 'peeler_theory_count', 2, 'INT', 'Теоретическое число чураков'),
(113, 12, 4, 'peeler_count', 2, 'INT', 'Число чураков за минуту'),
(114, 1, 10, 'pleer_sys_number', 4, 'INT', 'Номер чурака в системе'),
(115, 1, 11, 'pleer_last_number', 4, 'INT', 'Индекс последнего отправленного чурака'),
(116, 1, 12, 'service', 4, 'INT', 'Служебная информация'),
(117, 1, 13, 'knife_position', 4, 'INT', 'Открытие ножевой щели'),
(118, 13, 0, 'length', 2, 'INT', 'Length'),
(119, 13, 1, 'width', 2, 'INT', 'Width'),
(120, 13, 2, 'square', 2, 'INT', 'Square'),
(121, 13, 3, 'merges', 2, 'INT', 'Merges count');

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;

- N(程序段号)        G （准备功能字） X Y Z （坐标功能字）...  F（进给功能字） S（主轴转速功能字）T（刀具号）M（辅助功能）

  

- G01 走直线 G00快速点定位

- G02(顺时针)/G03(逆时针)走圆弧

- M30 程序结束指令。M98调用子程序。M99子程序结束

- M08开冷却，M09关冷却

- G92 螺纹单循环指令

  
  
  工序安排： 基准先行，先主后次，先粗后精，先面后孔



铣刀选择：

- 大平面： 面铣刀
- 加工凹槽，小台阶面及平面轮廓： 立铣刀
- 加工空间曲面，模具型腔等：球头铣刀
- 加工封闭的键槽：键槽铣刀等
- 加工变斜角零件： 鼓形铣刀
- 特殊形状： 成形铣刀
- 薄壁件加工，如整体框，肋类零件，为防止由于震动引起被动再切削，选用短切削刃的立铣刀

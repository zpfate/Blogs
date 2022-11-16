## 性能优化

1. widget的build方法中避免执行耗时操作

2. 控制setState的范围，可以使用多个StatefulWidget，Provider和Getx等控制刷新粒度的框架

3. 返回空widget的时候用nil代替Container和SizedBox

4. 控制widget的重建次数，使用StatelessWidget，const修饰，AnimatedBuilder将不需要动的子widget赋值给child参数

5. 使用RepaintBoundary

6. 耗时任务使用多线程isolate（快捷使用compute），future利用的是事件循环机制

   
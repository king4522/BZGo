#ifndef YOLOV5_H
#define YOLOV5_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct YoloV5 YoloV5;

YoloV5* newYoloV5();

const char *loadModelYoloV5(YoloV5 *obj, const char *param_path, const char *bin_path, const char* labels, int num_threads);

char *detectYoloV5(YoloV5 *obj, const unsigned char *bitmapData, int width, int height, float nms, float prob, int size);

void closeYoloV5(YoloV5 *obj);

#ifdef __cplusplus
}
#endif

#endif // YOLOV5_H

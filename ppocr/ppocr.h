#ifndef PPOCR_H
#define PPOCR_H

#ifdef __cplusplus
extern "C" {
#endif

typedef struct Ppocr Ppocr;

Ppocr* newPpocr();

const char* loadModelPpocr(Ppocr* Obj, const char* label, const char* dbParam, const char* dbBin, const char* recParam, const char* recBin, int thread);

char* detectPpocr(Ppocr* Obj, const unsigned char* bitmapData, int width, int height, float nms, float prob, int size);

void closePpocr(Ppocr* Obj);

#ifdef __cplusplus
}
#endif

#endif // PPOCR_H

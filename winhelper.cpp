#include "winhelper.h"

void freeUTF16String(LPWSTR str) {
	CoTaskMemFree(str);
}

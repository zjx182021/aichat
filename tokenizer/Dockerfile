FROM quay.io/0voice/python:3.10-alpine
WORKDIR /app
ENV PORT 3002
ADD ./tokenizer.py /app
ADD ./requirements.txt /app

RUN pip install -i https://pypi.tuna.tsinghua.edu.cn/simple --trusted-host pypi.tuna.tsinghua.edu.cn --upgrade pip
RUN pip install --root-user-action=ignore -i https://pypi.tuna.tsinghua.edu.cn/simple --trusted-host pypi.tuna.tsinghua.edu.cn -r requirements.txt

CMD ["sh", "-c", "nuxt --port ${PORT} --module tokenizer.py --workers 2"]

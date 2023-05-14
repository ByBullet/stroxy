import os
import shutil
import zipfile
import subprocess


def mkdir(path):
    folder = os.path.exists(path)
    if not folder:  # 判断是否存在文件夹如果不存在则创建为文件夹
        os.makedirs(path)  # makedirs 创建文件时如果路径不存在会创建这个路径
        print("---  new folder 【" + path + "】  ---")
    else:
        print("---  Folder 【" + path + "】 already exists!  ---")


def build(path, operate_system, plantform):
    mkdir(path)
    command = "build.bat {path} {operate_system} {plantform}"
    command = command.replace("{path}", path)
    command = command.replace("{operate_system}", operate_system)
    command = command.replace("{plantform}", plantform)
    print(command)
    os.system(command)


'''
编译golang程序
'''


def cross_compile_go_file(go_file_path, output_dir, target_os, target_arch, args):
    # 设置交叉编译时的 GOOS 和 GOARCH 环境变量
    os.environ["GOOS"] = target_os
    os.environ["GOARCH"] = target_arch

    output_filename = os.path.basename(go_file_path)
    if target_os == "windows":
        output_filename += ".exe"

    cmd = f"go build {args} -o {os.path.join(output_dir, output_filename)} {go_file_path} "
    subprocess.run(cmd, shell=True)


def dirCopy(source_path, target_path):
    if not os.path.exists(target_path):
        # 如果目标路径不存在原文件夹的话就创建
        os.makedirs(target_path)

    if os.path.exists(source_path):
        # 如果目标路径存在原文件夹的话就先删除
        shutil.rmtree(target_path)

    shutil.copytree(source_path, target_path)
    print('copy dir finished!')


def zip_file(src_dir):
    zip_name = src_dir + '.zip'
    z = zipfile.ZipFile(zip_name, 'w', zipfile.ZIP_DEFLATED)
    for dirpath, dirnames, filenames in os.walk(src_dir):
        fpath = dirpath.replace(src_dir, '')
        fpath = fpath and fpath + os.sep or ''
        for filename in filenames:
            z.write(os.path.join(dirpath, filename), fpath + filename)
            print('==压缩成功==')
    z.close()


if __name__ == "__main__":
    list = [
        # {"path": "win", "os": "windows", "plantform": "amd64", "args": '-ldflags "-H=windowsgui"'},
        # {"path": "mac", "os": "darwin", "plantform": "arm64","args": ""}
        # {"path":"linux","os":"linux","plantform":"amd64"},
        # {"path":"mac_arm64","os":"darwin","plantform":"arm64","args":""},
        # {"path":"mac_amd64","os":"darwin","plantform":"amd64","args":""}
    ]
    current_directory = os.path.dirname(os.path.abspath(__file__))
    path = current_directory + "/build/"
    if os.path.exists(path):
        # 如果目标路径存在原文件夹的话就先删除
        shutil.rmtree(path)

    for item in list:
        dst_dir = os.path.join(path, item["path"])
        shutil.copytree(current_directory + "/resources", dst_dir + "/resources")
        shutil.copytree(current_directory + "/script", dst_dir + "/script")
        cross_compile_go_file(os.path.join(current_directory, "app/gui"), dst_dir, item["os"], item["plantform"],
                              item["args"])
        # build(dst_dir,item["os"],item["plantform"])

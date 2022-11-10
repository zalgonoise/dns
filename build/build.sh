#!/bin/zsh
if [[ -z $1 ]]; then
    echo "no release version provided"
    exit 1
fi
if [[ -d ./build/out ]]; then
    rm -rf ./build/out
fi

mkdir -p ./build/out
touch ./build/out/README.md

release="$1"
out="./build/out"
package_name=$(pwd | sed 's/\// /g' | awk '{print $NF}')

os=(
    linux
    darwin
    windows
    android
    dragonfly
    freebsd
    netbsd
    openbsd
    plan9
    solaris
)

linux=(
    amd64
    386
    arm
    arm64
    ppc64
    ppc64le
    mips
    mipsle
    mips64
    mips64le
)

darwin=(
    amd64
    arm64
    # 386 # unsupported
    # arm # unsupported 
)

windows=( amd64 386 )

android=( arm64 )

dragonfly=( amd64 )

freebsd=( amd64 386 arm )

netbsd=( amd64 386 arm )

openbsd=( amd64 386 arm )

plan9=( amd64 386 )

solaris=( amd64 )

cat << EOF > ${out}/README.md
# ${package_name} - ${release}

*SHA256 Checksum list*

File | Sum
:---:|:----:
EOF

for (( a=1 ; a<=${#os[@]} ; a++ ))
do
    for (( b=1 ; b<=${(P)#os[a][@]} ; b++ ))
    do
        echo -n "Building for ${os[a]}_${(P)os[a][b]} -- "
        env GOOS=${os[a]} GOARCH=${(P)os[a][b]} go build -o ${out}/${package_name}_${os[a]}_${(P)os[a][b]} .
        if [[ -f ${out}/${package_name}_${os[a]}_${(P)os[a][b]} ]]
        then 
            echo "${package_name}_${os[a]}_${(P)os[a][b]} | \`$(sha256sum ${out}/${package_name}_${os[a]}_${(P)os[a][b]} | awk '{print $1}')\`" >> ${out}/README.md 
            echo "OK"

            if [[ ${os[a]} == "windows" ]]
            then
                mv ${out}/${package_name}_${os[a]}_${(P)os[a][b]} ${out}/${package_name}_${os[a]}_${(P)os[a][b]}.exe
            fi

        else 
            echo "Failed"
        fi

    done
done
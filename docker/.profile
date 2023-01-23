

alias ll='ls -laF'

export PATH=${PATH}:${HOME}/ezc/ezcluster/bin
export EZC_PLUGINS=${HOME}/ezc
export EZC_HELPERS=${HOME}/nih
export EZC_BUILDER_PRIVATE_KEY_PATH=${HOME}/.ezcluster
export EZC_CA_PATH=${HOME}/.ezcluster

if [ -z "$EZCT_LOG_MODE" ]
then
    export EZCT_LOG_MODE="dev"
fi

if [ ! -z "$EZCT_REPO" ]
then
    # We are in interactive mode
    # Will clone the repo
    tower refresh --localPath "."
    export BASE_DIR="$(tower location --localPath "." base)"
    export REPO_NAME="$(tower location --localPath "." reponame)"
    if [ "$REPO_NAME" == "osiac" ]
    then
        pattern='iac/*/*'
    fi
    if [ ! -z "$pattern" ]
    then
        cluster() {
          if [ -d $BASE_DIR/${pattern}/$1 ]
          then
            cd $BASE_DIR/${pattern}/$1
            if [ -x $BASE_DIR/${pattern}/$1/build/activate.sh ]
            then
              . $BASE_DIR/${pattern}/$1/build/activate.sh
            fi
            if [ -x $BASE_DIR/${pattern}/$1/.remote/activate.sh ]
            then
              . $BASE_DIR/${pattern}/$1/.remote/activate.sh
            fi
            return
          fi
          echo "Unable to locate $1"
        }
    fi
    cd $(tower location --localPath "." base)
else
    cd ${HOME}
fi



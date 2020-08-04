# This startup script creates a config file based on environment variables and then starts the bot

# CONFIG => holds the stormy.[json|yml|toml|hcl|env] content
# CONFIG_TYPE => json|yml|toml|hcl|env can be used, default: json
# CONFIG_2, CONFIG_3, CONFIG_4 => can be used to split the config into multiple variables
# ENCODED_CONFIG => set to true, to base64 decode the config file
# DEBUG => when set, the go binary will start in debug mode

config_filename=stormy.${CONFIG_TYPE:-json}

# read config, split in up to 4 environment variables, and create the config file
echo ${CONFIG}${CONFIG_2}${CONFIG_3}${CONFIG_4} > $config_filename

# decode the config file, if it's base64 encoded
if [ -n ${ENCODED_CONFIG} ] ; 
  then base64 -d $config_filename | tee $config_filename; 
  else cat $config_filename ;
fi

# start the bot
if [ -n ${DEBUG} ] ; 
  then ./bin/stormy -debug; 
  else ./bin/stormy ; 
fi

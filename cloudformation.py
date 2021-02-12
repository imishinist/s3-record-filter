#!/usr/bin/env python
# -*- coding: utf-8 -*-

import sys
import os.path

import boto3
import json
import _jsonnet
import fire

from awscli.customizations.cloudformation import exceptions
from awscli.customizations.cloudformation.deployer import Deployer

from awscli.compat import get_stdout_text_writer
from awscli.utils import write_exception


class CFnDeploy(object):
    _MSG_NO_EXECUTE_CHANGESET = \
        ("Changeset created successfully. Run the following command to "
         "review changes:"
         "\n"
         "aws cloudformation describe-change-set --change-set-name "
         "{changeset_id}"
         "\n")
    _MSG_EXECUTE_SUCCESS = "Successfully created/updated stack - {stack_name}\n"

    def __load_jsonnet(self, params_path):
        with open(params_path) as f:
            jsonnet_str = f.read()

        json_str = _jsonnet.evaluate_snippet(
            "snippet", jsonnet_str)

        return json.loads(json_str)

    def __load_json(self, params_path):
        with open(params_path) as f:
            return json.load(f)

    def __load_cfn_params(self, params_path):
        ext = os.path.splitext(params_path)[1]
        if ext == ".jsonnet":
            params = self.__load_jsonnet(params_path)
        else:
            params = self.__load_json(params_path)

        cfn_params = []
        for key in params:
            cfn_params.append(
                {
                    'ParameterKey': str(key),
                    'ParameterValue': str(params[key]),
                    'UsePreviousValue': False,
                }
            )
        return cfn_params

    def __deploy(self, deployer, stack_name, template_str,
                 parameters, capabilities, execute_changeset, role_arn,
                 notification_arns, s3_uploader, tags, fail_on_empty_changeset=True):
        try:
            result = deployer.create_and_wait_for_changeset(
                stack_name=stack_name,
                cfn_template=template_str,
                parameter_values=parameters,
                capabilities=capabilities,
                role_arn=role_arn,
                notification_arns=notification_arns,
                s3_uploader=s3_uploader,
                tags=tags
            )
        except exceptions.ChangeEmptyError as ex:
            if fail_on_empty_changeset:
                raise
            write_exception(ex, outfile=get_stdout_text_writer())
            return 0

        if execute_changeset:
            deployer.execute_changeset(result.changeset_id, stack_name)
            deployer.wait_for_execute(stack_name, result.changeset_type)
            sys.stdout.write(self._MSG_EXECUTE_SUCCESS.format(stack_name=stack_name))
        else:
            sys.stdout.write(self.MSG_NO_EXECUTE_SUCCESS.format(changeset_id=result.changeset_id))

        sys.stdout.flush()
        return 0

    def __init__(self, stack_name, template_path, params_path):
        self.stack_name = stack_name
        self.template_path = template_path
        self.params_path = params_path

    # ported from:
    # https://github.com/aws/aws-cli/blob/develop/awscli/customizations/cloudformation/deploy.py
    def deploy(self):
        """
        [option]

        option:
        stack-name: CloudFormation stack name
        template-path: CloudFormation template file path
        params-path: CloudFormation template parameters file path.
                     allowed file type: json, jsonnet
        """
        client = boto3.client('cloudformation')

        if not os.path.isfile(self.template_path):
            raise exceptions.InvalidTemplatePathError(template_path=self.template_path)
        with open(self.template_path, "r") as f:
            template_str = f.read()

        params = self.__load_cfn_params(self.params_path)
        tags = []
        s3_uploader = None

        deployer = Deployer(client)
        self.__deploy(deployer, self.stack_name, template_str,
                      params, ['CAPABILITY_IAM', 'CAPABILITY_AUTO_EXPAND'],
                      True, None, None, s3_uploader, tags, False)

if __name__ == '__main__':
    fire.Fire(CFnDeploy)

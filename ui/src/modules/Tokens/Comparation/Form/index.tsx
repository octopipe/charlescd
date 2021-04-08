/*
 * Copyright 2020 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { Fragment } from 'react';
import { FormProvider, useForm } from 'react-hook-form';
import { Token } from 'modules/Tokens/interfaces';
import { useSave } from 'modules/Tokens/hooks';
import ContentIcon from 'core/components/ContentIcon';
import Form from 'core/components/Form';
import { isRequiredAndNotBlank } from 'core/utils/validations';
import Workspaces from './Workspaces';
import Scopes from './Scopes';
import Styled from './styled';
import { Mode } from '../helpers';

interface Props {
  mode?: Mode;
}

const FormToken = ({ mode }: Props) => {
  const { save, status } = useSave();

  const methods = useForm<Token>({ mode: 'onChange' });
  const {
    register, handleSubmit, watch,
    setValue, errors, formState: { isValid }
  } = methods;

  const name = watch('name') as string;
  const workspaces = watch('workspaces') as string[];

  const onSubmit = (token: Token) => {
    const { subjects, ...rest } = token;
    save(rest);
  };

  return (
    <Styled.Content>
      <FormProvider {...methods}>
        <Styled.Form onSubmit={handleSubmit(onSubmit)}>
          <ContentIcon icon="token">
            <Form.InputTitle
              name="name"
              ref={register(isRequiredAndNotBlank)}
              error={errors?.name?.message}
            />
          </ContentIcon>
          {name && <Workspaces setValue={setValue} />}
          {name && workspaces && (
            <Fragment>
              <Scopes mode={mode} />
              <Styled.Button
                type="submit"
                size="EXTRA_SMALL"
                isDisabled={!isValid}
                isLoading={status.isPending}
              >
                Generate token
              </Styled.Button>
            </Fragment>
          )}
        </Styled.Form>
      </FormProvider>
    </Styled.Content>
  );
};

export default FormToken;

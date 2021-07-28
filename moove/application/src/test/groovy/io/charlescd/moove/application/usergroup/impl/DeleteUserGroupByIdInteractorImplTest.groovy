/*
 * Copyright 2020, 2021 ZUP IT SERVICOS EM TECNOLOGIA E INOVACAO SA
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

package io.charlescd.moove.application.usergroup.impl

import io.charlescd.moove.application.UserGroupService
import io.charlescd.moove.application.usergroup.DeleteUserGroupByIdInteractor
import io.charlescd.moove.domain.User
import io.charlescd.moove.domain.UserGroup
import io.charlescd.moove.domain.exceptions.NotFoundException
import io.charlescd.moove.domain.repository.UserGroupRepository
import spock.lang.Specification

import java.time.LocalDateTime

class DeleteUserGroupByIdInteractorImplTest extends Specification {

    private DeleteUserGroupByIdInteractor deleteUserGroupByIdInteractor

    private UserGroupRepository userGroupRepository = Mock(UserGroupRepository)

    void setup() {
        this.deleteUserGroupByIdInteractor = new DeleteUserGroupByIdInteractorImpl(new UserGroupService(userGroupRepository))
    }

    def "when user group does not exists should throw an exception"() {
        given:
        def userGroupId = "user-group-id"

        when:
        this.deleteUserGroupByIdInteractor.execute(userGroupId)

        then:
        1 * this.userGroupRepository.findById(userGroupId) >> Optional.empty()

        def ex = thrown(NotFoundException)
        ex.resourceName == "user_group"
        ex.id == userGroupId
    }

    def "should find and delete user group by id"() {
        given:
        def authorId = "0a859e6c-3cdf-4b34-84d0-f9038576ac58"
        def author = new User(authorId, "charles", "charles@zup.com.br", "http://charles.com/dummy_photo.jpg", [], [], false, LocalDateTime.now())
        def userGroupId = "user-group-id"
        def userGroup = new UserGroup(userGroupId, "group-name", author, LocalDateTime.now(), [])

        when:
        this.deleteUserGroupByIdInteractor.execute(userGroupId)

        then:
        1 * this.userGroupRepository.findById(userGroupId) >> Optional.of(userGroup)
        1 * this.userGroupRepository.delete(userGroup) >> _

        notThrown()
    }
}

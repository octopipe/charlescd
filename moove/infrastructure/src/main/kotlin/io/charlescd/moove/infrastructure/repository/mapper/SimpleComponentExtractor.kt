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

package io.charlescd.moove.infrastructure.repository.mapper

import io.charlescd.moove.domain.SimpleComponent
import java.sql.ResultSet
import org.springframework.jdbc.core.ResultSetExtractor

@org.springframework.stereotype.Component
class SimpleComponentExtractor : ResultSetExtractor<Set<SimpleComponent>> {

    override fun extractData(resultSet: ResultSet): Set<SimpleComponent>? {
        val componentsResult = HashSet<SimpleComponent>()

        while (resultSet.next()) {
            componentsResult.add(mapSimpleComponent(resultSet))
        }

        return componentsResult
    }

    private fun mapSimpleComponent(resultSet: ResultSet) = SimpleComponent(
            id = resultSet.getString("components_id"),
            moduleId = resultSet.getString("components_module_id"),
            name = resultSet.getString("components_name"),
            errorThreshold = resultSet.getInt("components_error_threshold"),
            latencyThreshold = resultSet.getInt("components_latency_threshold"),
            workspaceId = resultSet.getString("components_workspace_id"),
            moduleName = resultSet.getString("module_name")
    )
}
